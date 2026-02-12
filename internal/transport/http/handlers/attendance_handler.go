package handlers

import (
	"autoattendance-go/internal/domain"
	"autoattendance-go/internal/dto"
	"autoattendance-go/internal/service"
	"net/http"

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
	role := c.MustGet("user_role").(domain.Role)
	if role == domain.RoleEmployee {
		req.UserID = userID
	} else if req.UserID == uuid.Nil {
		req.UserID = userID
	}

	req.AgencyID = agencyID

	res, err := h.svc.MarkAttendance(c.Request.Context(), &req)
	if err != nil {
		if err == domain.ErrAttendanceExists {
			c.JSON(http.StatusConflict, gin.H{"error": "asistencia ya registrada para hoy"})
			return
		}
		if err == domain.ErrNoScheduleFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no tienes un horario asignado para hoy"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
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
	role := c.MustGet("user_role").(domain.Role)
	if role == domain.RoleEmployee {
		params.UserID = c.MustGet("user_id").(uuid.UUID)
	}

	res, err := h.svc.GetAgencyAttendances(c.Request.Context(), agencyID, &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}

	c.JSON(http.StatusOK, res)
}
