package handlers

import (
	"net/http"
	"quickattendance-go/internal/dto"
	"quickattendance-go/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ScheduleHandler struct {
	svc *service.ScheduleService
}

func NewScheduleHandler(svc *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{svc: svc}
}

func (h *ScheduleHandler) Create(c *gin.Context) {
	agencyID := c.MustGet("agency_id").(uuid.UUID)
	var req dto.CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.svc.CreateSchedule(c.Request.Context(), agencyID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *ScheduleHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	scheduleID, err := uuid.Parse(idStr)
	if err != nil || scheduleID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule ID"})
		return
	}

	agencyID := c.MustGet("agency_id").(uuid.UUID)

	res, err := h.svc.GetSchedule(c.Request.Context(), scheduleID, agencyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "schedule not found"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *ScheduleHandler) List(c *gin.Context) {
	agencyID := c.MustGet("agency_id").(uuid.UUID)

	var params dto.ScheduleListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Now passing params for filtering and pagination
	res, err := h.svc.GetAgencySchedules(c.Request.Context(), agencyID, &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *ScheduleHandler) GetApplicable(c *gin.Context) {
	agencyID := c.MustGet("agency_id").(uuid.UUID)

	var params dto.GetApplicableScheduleParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)
	// Simple role check for demo (could use middleware)
	// role := c.MustGet("user_role").(string)
	if params.UserID != uuid.Nil {
		userID = params.UserID
	}

	parsedDate, err := time.Parse("2006-01-02", params.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
		return
	}

	res, err := h.svc.GetApplicableSchedule(c.Request.Context(), agencyID, userID, parsedDate)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no applicable schedule found"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *ScheduleHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	scheduleID, err := uuid.Parse(idStr)
	if err != nil || scheduleID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule ID"})
		return
	}

	agencyID := c.MustGet("agency_id").(uuid.UUID)

	var req dto.UpdateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.svc.UpdateSchedule(c.Request.Context(), &req, scheduleID, agencyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *ScheduleHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	scheduleID, err := uuid.Parse(idStr)
	if err != nil || scheduleID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule ID"})
		return
	}

	agencyID := c.MustGet("agency_id").(uuid.UUID)

	err = h.svc.DeleteSchedule(c.Request.Context(), scheduleID, agencyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "schedule deleted"})
}
