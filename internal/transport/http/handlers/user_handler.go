package handlers

import (
	"net/http"
	"quickattendance-go/internal/domain"
	"quickattendance-go/internal/dto"
	"quickattendance-go/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// Invite godoc
// @Summary Invite a new user to the agency
// @Description Creates a new user record and sends an invitation email (Admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.InviteUserRequest true "Invitation details"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /users/invite [post]
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

// Activate godoc
// @Summary Activate a user account
// @Description Activates a user account using the code sent via email. Returns JWT token.
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.ActivateUserRequest true "Activation code and credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/activate [post]
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

// Login godoc
// @Summary Login user
// @Description Authenticate user with email and password. Returns JWT token.
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.LoginUserRequest true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/login [post]
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

// GetMe godoc
// @Summary Get current user profile
// @Description Returns the profile of the currently authenticated user.
// @Tags users
// @Produce json
// @Success 200 {object} domain.User
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /users/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	agencyID := c.MustGet("agency_id").(uuid.UUID)
	if userID == uuid.Nil || agencyID == uuid.Nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	user, err := h.svc.GetUserByID(c.Request.Context(), agencyID, userID)

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

func (h *UserHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := uuid.Parse(idStr)
	if err != nil || userID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	agencyID := c.MustGet("agency_id").(uuid.UUID)

	user, err := h.svc.GetUserByID(c.Request.Context(), agencyID, userID)
	if err != nil {
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
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
	if err != nil || userID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is invalid"})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.svc.UpdateUserProfile(c.Request.Context(), agencyID, userID, &req)
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
	if err != nil || userID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}
	agencyID := c.MustGet("agency_id").(uuid.UUID)

	err = h.svc.DeleteUser(c.Request.Context(), agencyID, userID)
	if err != nil {
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
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

	var params dto.UserListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Now passing params for filtering and pagination
	users, err := h.svc.ListByAgencyID(c.Request.Context(), agencyID, &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, users)
}
