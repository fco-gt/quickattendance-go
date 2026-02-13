package handlers

import (
	"net/http"
	"quickattendance-go/internal/domain"
	"quickattendance-go/internal/dto"
	"quickattendance-go/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AttendanceHandler struct {
	svc *service.AttendanceService
}

func NewAttendanceHandler(svc *service.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{svc: svc}
}

func (h *AttendanceHandler) Mark(c *gin.Context) {
	agencyID := c.MustGet("agency_id").(uuid.UUID)
	userID := c.MustGet("user_id").(uuid.UUID)

	var req dto.MarkAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Security: Employees can only mark their own attendance
	role := c.MustGet("role").(domain.Role)
	req.RequesterRole = role
	if role == domain.RoleEmployee {
		req.UserID = userID
	} else if req.UserID == uuid.Nil {
		req.UserID = userID
	}

	req.AgencyID = agencyID

	res, err := h.svc.MarkAttendance(c.Request.Context(), &req)
	if err != nil {
		if err == domain.ErrAttendanceExists {
			c.JSON(http.StatusConflict, gin.H{"error": "attendance already registered for today"})
			return
		}
		if err == domain.ErrNoScheduleFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no schedule found for today"})
			return
		}
		if err == domain.ErrManualNotAllowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "only admins can mark attendance manually"})
			return
		}
		if err == domain.ErrGeofenceViolation {
			c.JSON(http.StatusBadRequest, gin.H{"error": "you are out of the allowed range from your home"})
			return
		}
		if err == domain.ErrHomeLocationNotSet {
			c.JSON(http.StatusBadRequest, gin.H{"error": "home location not configured for remote marking"})
			return
		}
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if err == domain.ErrInvalidAttendance {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attendance data"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *AttendanceHandler) List(c *gin.Context) {
	agencyID := c.MustGet("agency_id").(uuid.UUID)

	var params dto.AttendanceListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Security: Employees can only see their own attendance
	role := c.MustGet("role").(domain.Role)
	if role == domain.RoleEmployee {
		params.UserID = c.MustGet("user_id").(uuid.UUID).String()
	} else if params.UserID != "" {
		if _, err := uuid.Parse(params.UserID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
			return
		}
	}

	res, err := h.svc.GetAgencyAttendances(c.Request.Context(), agencyID, &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, res)
}
