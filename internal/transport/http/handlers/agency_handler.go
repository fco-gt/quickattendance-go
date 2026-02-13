package handlers

import (
	"net/http"
	"quickattendance-go/internal/domain"
	"quickattendance-go/internal/dto"
	"quickattendance-go/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AgencyHandler struct {
	svc *service.AgencyService
}

func NewAgencyHandler(svc *service.AgencyService) *AgencyHandler {
	return &AgencyHandler{svc: svc}
}

// Register godoc
// @Summary Register a new agency
// @Description Creates a new agency and its first admin user.
// @Tags agencies
// @Accept json
// @Produce json
// @Param request body dto.RegisterAgencyRequest true "Agency registration details"
// @Success 201 {object} dto.AgencyResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /agencies [post]
func (h *AgencyHandler) Register(c *gin.Context) {
	var req dto.RegisterAgencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data: " + err.Error()})
		return
	}

	res, err := h.svc.Register(c.Request.Context(), &req)

	if err != nil {
		if err == domain.ErrAgencyExists || err == domain.ErrUserExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, res)
}

// Update godoc
// @Summary Update agency details
// @Description Updates the details of the agency (Admin only).
// @Tags agencies
// @Accept json
// @Produce json
// @Param request body dto.UpdateAgencyRequest true "Agency updated details"
// @Success 200 {object} domain.Agency
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /agencies [put]
func (h *AgencyHandler) Update(c *gin.Context) {
	agencyID := c.MustGet("agency_id").(uuid.UUID)
	if agencyID == uuid.Nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	var req dto.UpdateAgencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data: " + err.Error()})
		return
	}

	agency, err := h.svc.Update(c.Request.Context(), agencyID, &req)
	if err != nil {
		if err == domain.ErrAgencyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, agency)
}
