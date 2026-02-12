package domain

import (
	"context"
)

type NotificationProvider interface {
	PublishEmail(ctx context.Context, to string, subject string, body string) error
}
