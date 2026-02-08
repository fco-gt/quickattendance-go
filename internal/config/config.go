package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Env            string
	HTTPPort       string
	DatabaseURL    string
	JWTSecret      string
	AccessTokenTTL time.Duration
	BCryptCost     int
}

func Load() *Config {
	ttlStr := getEnv("ACCESS_TOKEN_TTL", "15m")
	bcryptCostStr := getEnv("BCRYPT_COST", "12")

	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		log.Fatalf("invalid ACCESS_TOKEN_TTL: %v", err)
	}

	bcryptCost, err := strconv.Atoi(bcryptCostStr)
	if err != nil {
		log.Fatalf("invalid BCRYPT_COST: %v", err)
	}

	cfg := &Config{
		Env:            getEnv("APP_ENV", "development"),
		HTTPPort:       getEnv("HTTP_PORT", "8080"),
		DatabaseURL:    mustEnv("DB_URL"),
		JWTSecret:      mustEnv("JWT_SECRET"),
		AccessTokenTTL: ttl,
		BCryptCost:     bcryptCost,
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

func mustEnv(key string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	log.Fatalf("missing required env var: %s", key)
	return ""
}
