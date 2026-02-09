package main

import (
	"autoattendance-go/internal/config"
	"autoattendance-go/internal/domain"
	"autoattendance-go/internal/repository"
	"autoattendance-go/internal/service"
	"autoattendance-go/internal/transport/http/handlers"
	"autoattendance-go/pkg/logger"
	"autoattendance-go/pkg/security"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Config
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found")
	}
	cfg := config.Load()

	// Logger Setup
	logger.Setup(cfg.Env)

	// Database with retries (Docker might take time so we need to wait)
	var db *gorm.DB
	var err error
	for i := range 10 {
		db, err = gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
		if err == nil {
			break
		}
		slog.Info("Waiting for database...", "attempt", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	db.AutoMigrate(&domain.Agency{}, &domain.User{})

	// Utilities
	jwtService := security.NewJWTService(cfg.JWTSecret)
	hasher := security.NewPasswordHasher(cfg.BCryptCost)
	tokenTTL := cfg.AccessTokenTTL

	// Repositories
	agencyRepo := repository.NewAgencyRepo(db)
	userRepo := repository.NewUserRepo(db)
	txManager := repository.NewGormTransactor(db)

	// Services
	agencySvc := service.NewAgencyService(agencyRepo, userRepo, hasher, txManager)
	userSvc := service.NewUserService(userRepo, agencyRepo, jwtService, hasher, tokenTTL)

	// Router
	r := handlers.NewRouter(agencySvc, userSvc, jwtService)

	// Server
	fmt.Printf("Server running on port %s\n", cfg.HTTPPort)
	if err := r.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatal("failed to run server:", err)
	}
}
