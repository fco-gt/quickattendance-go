package security

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type jobType int

const (
	jobHash jobType = iota
	jobCompare
)

type job struct {
	kind     jobType
	password string
	hash     string
	result   chan result
}

type result struct {
	hash string
	ok   bool
	err  error
}

type PasswordHasher struct {
	jobs    chan job
	workers int
}

func NewPasswordHasher(workers int) *PasswordHasher {
	h := &PasswordHasher{
		jobs:    make(chan job, 100),
		workers: workers,
	}

	for range workers {
		go h.worker()
	}

	return h
}

func (h *PasswordHasher) worker() {
	for j := range h.jobs {
		switch j.kind {
		case jobHash:
			hash, err := bcrypt.GenerateFromPassword(
				[]byte(j.password),
				bcrypt.DefaultCost,
			)
			j.result <- result{
				hash: string(hash),
				err:  err,
			}
		case jobCompare:
			err := bcrypt.CompareHashAndPassword(
				[]byte(j.hash),
				[]byte(j.password),
			)
			j.result <- result{
				ok:  err == nil,
				err: err,
			}
		}
	}
}

func (h *PasswordHasher) Hash(ctx context.Context, password string) (string, error) {

	res := make(chan result, 1)

	select {
	case h.jobs <- job{
		kind:     jobHash,
		password: password,
		result:   res,
	}:
		select {
		case r := <-res:
			return r.hash, r.err
		case <-ctx.Done():
			return "", ctx.Err()
		}

	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func (h *PasswordHasher) Compare(
	ctx context.Context,
	hash string,
	password string,
) (bool, error) {

	res := make(chan result, 1)

	select {
	case h.jobs <- job{
		kind:     jobCompare,
		hash:     hash,
		password: password,
		result:   res,
	}:
		select {
		case r := <-res:
			return r.ok, r.err
		case <-ctx.Done():
			return false, ctx.Err()
		}

	case <-ctx.Done():
		return false, ctx.Err()
	}
}
