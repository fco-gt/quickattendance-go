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
