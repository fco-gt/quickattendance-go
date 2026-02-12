package handlers

import (
	"quickattendance-go/internal/domain"
	"quickattendance-go/internal/service"
	"quickattendance-go/internal/transport/http/middleware"
	"quickattendance-go/pkg/security"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	agencySvc *service.AgencyService,
	userSvc *service.UserService,
	scheduleSvc *service.ScheduleService,
	attendanceSvc *service.AttendanceService,
	jwtSvc *security.JWTService,
) *gin.Engine {
	r := gin.Default()

	// Handlers
	agencyHandler := NewAgencyHandler(agencySvc)
	userHandler := NewUserHandler(userSvc)
	scheduleHandler := NewScheduleHandler(scheduleSvc)
	attendanceHandler := NewAttendanceHandler(attendanceSvc)

	// Middlewares
	authMiddleware := middleware.Auth(jwtSvc)

	// Basic CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	v1 := r.Group("/api/v1")
	{
		// Health
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"service": "quickattendance-api",
			})
		})

		// Agencies routes
		agencies := v1.Group("/agencies")
		{
			agencies.POST("", agencyHandler.Register)

			// Protected routes
			protected := agencies.Group("")
			protected.Use(authMiddleware)
			{
				protected.PUT("", middleware.RequireRole(domain.RoleAdmin), agencyHandler.Update)
			}
		}

		// Users routes
		users := v1.Group("users")
		{
			users.POST("/activate", userHandler.Activate)
			users.POST("/login", userHandler.Login)

			// Protected routes
			protected := users.Group("")
			protected.Use(authMiddleware)
			{
				protected.GET("/me", userHandler.GetMe)
				protected.POST("/invite", middleware.RequireRole(domain.RoleAdmin), userHandler.Invite)
				protected.GET("/:id", userHandler.GetByID)
				protected.PUT("/:id", middleware.RequireRole(domain.RoleAdmin), userHandler.UpdateProfile)
				protected.DELETE("/:id", middleware.RequireRole(domain.RoleAdmin), userHandler.Delete)
				protected.GET("/list", middleware.RequireRole(domain.RoleAdmin), userHandler.List)
			}
		}

		// Schedules routes
		schedules := v1.Group("schedules")
		schedules.Use(authMiddleware)
		{
			schedules.GET("/applicable", scheduleHandler.GetApplicable)
			schedules.GET("/list", scheduleHandler.List)
			schedules.GET("/:id", scheduleHandler.GetByID)

			// Admin only
			adminOnly := schedules.Group("")
			adminOnly.Use(middleware.RequireRole(domain.RoleAdmin))
			{
				adminOnly.POST("", scheduleHandler.Create)
				adminOnly.PUT("/:id", scheduleHandler.Update)
				adminOnly.DELETE("/:id", scheduleHandler.Delete)
			}
		}

		// Attendance routes
		attendance := v1.Group("attendance")
		attendance.Use(authMiddleware)
		{
			attendance.POST("/mark", attendanceHandler.Mark)
			attendance.GET("/list", attendanceHandler.List)
		}
	}
	return r
}
