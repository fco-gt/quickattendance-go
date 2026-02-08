package handlers

import (
	"autoattendance-go/internal/domain"
	"autoattendance-go/internal/dto"
	"autoattendance-go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Invite(c *gin.Context) {
	agencyId := c.MustGet("agency_id").(uuid.UUID)
	if agencyId == uuid.Nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	var req dto.InviteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.Invite(c.Request.Context(), agencyId, &req); err != nil {
		if err == domain.ErrUserExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user invited successfully"})
}

func (h *UserHandler) Activate(c *gin.Context) {
	var req dto.ActivateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.svc.ActivateByCode(c.Request.Context(), &req)
	if err != nil {
		if err == domain.ErrInvalidActivationCode || err == domain.ErrActivationCodeExpired {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err == domain.ErrUserAlreadyActivated {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": res.Token, "user": res.User})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.svc.Login(c.Request.Context(), &req)
	if err != nil {
		if err == domain.ErrInvalidCredentials || err == domain.ErrUserNotActive {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": res.Token, "user": res.User})
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	if userID == uuid.Nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	user, err := h.svc.GetUserByID(c.Request.Context(), userID)

	if err != nil {
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	idStr := c.Param("id")
	agencyID := c.MustGet("agency_id").(uuid.UUID)

	userID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is invalid"})
		return
	}

	if userID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is invalid"})
		return
	}

	user, err := h.svc.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if agencyID != user.AgencyID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err = h.svc.UpdateUserProfile(c.Request.Context(), userID, &req)
	if err != nil {
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}
	agencyID := c.MustGet("agency_id").(uuid.UUID)

	if userID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.svc.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if agencyID != user.AgencyID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	err = h.svc.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func (h *UserHandler) List(c *gin.Context) {
	agencyID := c.MustGet("agency_id").(uuid.UUID)

	if agencyID == uuid.Nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	users, err := h.svc.ListByAgencyID(c.Request.Context(), agencyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, users)
}
