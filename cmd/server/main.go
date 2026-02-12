package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"quickattendance-go/internal/config"
	"quickattendance-go/internal/domain"
	"quickattendance-go/internal/repository"
	"quickattendance-go/internal/service"
	"quickattendance-go/internal/transport/http/handlers"
	"quickattendance-go/pkg/logger"
	"quickattendance-go/pkg/messaging"
	"quickattendance-go/pkg/security"
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

	db.AutoMigrate(&domain.Agency{}, &domain.User{}, &domain.Schedule{}, &domain.Attendance{})

	// Utilities
	jwtService := security.NewJWTService(cfg.JWTSecret)
	hasher := security.NewPasswordHasher(cfg.BCryptCost)
	tokenTTL := cfg.AccessTokenTTL

	// RabbitMQ
	emailProducer, err := messaging.NewRabbitMQProducer(cfg.RabbitURL, "email_queue")
	if err != nil {
		slog.Error("failed to connect to RabbitMQ", "error", err)
		os.Exit(1)
	}
	defer emailProducer.Close()

	// Repositories
	agencyRepo := repository.NewAgencyRepo(db)
	userRepo := repository.NewUserRepo(db)
	scheduleRepo := repository.NewScheduleRepo(db)
	attendanceRepo := repository.NewAttendanceRepo(db)
	txManager := repository.NewGormTransactor(db)

	// Services
	agencySvc := service.NewAgencyService(agencyRepo, userRepo, hasher, txManager)
	userSvc := service.NewUserService(userRepo, agencyRepo, jwtService, hasher, emailProducer, tokenTTL)
	scheduleSvc := service.NewScheduleService(scheduleRepo, userRepo, txManager)
	attendanceSvc := service.NewAttendanceService(attendanceRepo, userRepo, scheduleSvc, txManager)

	// Router
	r := handlers.NewRouter(agencySvc, userSvc, scheduleSvc, attendanceSvc, jwtService)

	// Server
	fmt.Printf("Server running on port %s\n", cfg.HTTPPort)
	if err := r.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatal("failed to run server:", err)
	}
}
