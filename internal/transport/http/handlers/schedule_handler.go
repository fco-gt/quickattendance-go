package handlers

import (
	"net/http"
	"quickattendance-go/internal/domain"
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

// Create godoc
// @Summary Create a new schedule
// @Description Creates a new work schedule for the agency (Admin only).
// @Tags schedules
// @Accept json
// @Produce json
// @Param request body dto.CreateScheduleRequest true "Schedule details"
// @Success 201 {object} domain.Schedule
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /schedules [post]
func (h *ScheduleHandler) Create(c *gin.Context) {
	agencyID := c.MustGet("agency_id").(uuid.UUID)
	var req dto.CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.svc.CreateSchedule(c.Request.Context(), agencyID, &req)
	if err != nil {
		switch err {
		case domain.ErrDefaultScheduleAlreadyExists, domain.ErrScheduleNameAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, res)
}

// GetByID godoc
// @Summary Get schedule by ID
// @Description Returns a specific schedule details.
// @Tags schedules
// @Produce json
// @Param id path string true "Schedule ID"
// @Success 200 {object} domain.Schedule
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /schedules/{id} [get]
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

// List godoc
// @Summary List all schedules
// @Description Returns a list of schedules for the agency.
// @Tags schedules
// @Produce json
// @Success 200 {array} domain.Schedule
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /schedules/list [get]
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list schedules"})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetApplicable godoc
// @Summary Get applicable schedule for a date
// @Description Returns the schedule that applies to a user on a specific date.
// @Tags schedules
// @Produce json
// @Param user_id query string false "User ID (defaults to current user)"
// @Param date query string true "Date (YYYY-MM-DD)"
// @Success 200 {object} domain.Schedule
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /schedules/applicable [get]
func (h *ScheduleHandler) GetApplicable(c *gin.Context) {
	agencyID := c.MustGet("agency_id").(uuid.UUID)

	var params dto.GetApplicableScheduleParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)
	if params.UserID != "" {
		parsedID, err := uuid.Parse(params.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
			return
		}
		userID = parsedID
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

// Update godoc
// @Summary Update a schedule
// @Description Updates an existing schedule (Admin only).
// @Tags schedules
// @Accept json
// @Produce json
// @Param id path string true "Schedule ID"
// @Param request body dto.UpdateScheduleRequest true "Updated details"
// @Success 200 {object} domain.Schedule
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Security BearerAuth
// @Router /schedules/{id} [put]
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
		switch err {
		case domain.ErrScheduleNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case domain.ErrScheduleNameAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, res)
}

// Delete godoc
// @Summary Delete a schedule
// @Description Deletes a schedule from the agency (Admin only).
// @Tags schedules
// @Produce json
// @Param id path string true "Schedule ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /schedules/{id} [delete]
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
		switch err {
		case domain.ErrScheduleNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case domain.ErrDeleteDefaultSchedule:
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "schedule deleted"})
}
