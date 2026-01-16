package handlers

import (
	"errors"

	_ "github.com/carllix/matchaciee-backend/docs"
	"github.com/carllix/matchaciee-backend/internal/repositories"
	"github.com/carllix/matchaciee-backend/internal/services"
	"github.com/carllix/matchaciee-backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user account with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body docs.RegisterRequest true "Registration details"
// @Success 201 {object} docs.AuthSuccessResponse "User registered successfully"
// @Failure 400 {object} docs.SwaggerValidationErrorResponse "Validation error"
// @Failure 409 {object} docs.SwaggerErrorResponse "Email already exists"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req services.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return utils.ValidationErrorResponse(c, validationErrors)
	}

	// Register user
	authResp, err := h.authService.Register(req)
	if err != nil {
		if errors.Is(err, repositories.ErrEmailAlreadyExists) {
			return utils.ErrorResponse(c, fiber.StatusConflict, "Email already exists")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to register user")
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, authResp)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body docs.LoginRequest true "Login credentials"
// @Success 200 {object} docs.AuthSuccessResponse "Login successful"
// @Failure 400 {object} docs.SwaggerValidationErrorResponse "Validation error"
// @Failure 401 {object} docs.SwaggerErrorResponse "Invalid email or password"
// @Failure 403 {object} docs.SwaggerErrorResponse "User account is inactive"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req services.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return utils.ValidationErrorResponse(c, validationErrors)
	}

	// Login user
	authResp, err := h.authService.Login(req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid email or password")
		}
		if errors.Is(err, services.ErrUserInactive) {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "User account is inactive")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to login")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, authResp)
}

// GetMe godoc
// @Summary Get current user profile
// @Description Get the authenticated user's profile information
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} docs.MeSuccessResponse "User profile retrieved successfully"
// @Failure 401 {object} docs.SwaggerErrorResponse "Unauthorized"
// @Failure 404 {object} docs.SwaggerErrorResponse "User not found"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /auth/me [get]
func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	// Get user UUID from context
	userUUID := c.Locals("userUUID")
	if userUUID == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	uuid, ok := userUUID.(uuid.UUID)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Invalid user UUID")
	}

	// Get user
	user, err := h.authService.GetUserByUUID(uuid)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get user")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"user": user,
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get a new access token using a valid refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body docs.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} docs.AuthSuccessResponse "Token refreshed successfully"
// @Failure 400 {object} docs.SwaggerValidationErrorResponse "Validation error"
// @Failure 401 {object} docs.SwaggerErrorResponse "Invalid or expired refresh token"
// @Failure 403 {object} docs.SwaggerErrorResponse "User account is inactive"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req services.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return utils.ValidationErrorResponse(c, validationErrors)
	}

	// Refresh token
	authResp, err := h.authService.RefreshToken(req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidRefreshToken) {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired refresh token")
		}
		if errors.Is(err, services.ErrUserInactive) {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "User account is inactive")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to refresh token")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, authResp)
}

// Logout godoc
// @Summary Logout user
// @Description Invalidate the refresh token to logout the user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body docs.LogoutRequest true "Refresh token to invalidate"
// @Success 200 {object} docs.MessageSuccessResponse "Logged out successfully"
// @Failure 400 {object} docs.SwaggerValidationErrorResponse "Validation error"
// @Failure 401 {object} docs.SwaggerErrorResponse "Invalid refresh token"
// @Failure 500 {object} docs.SwaggerErrorResponse "Internal server error"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return utils.ValidationErrorResponse(c, validationErrors)
	}

	// Logout user
	err := h.authService.Logout(req.RefreshToken)
	if err != nil {
		if errors.Is(err, services.ErrInvalidRefreshToken) {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid refresh token")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to logout")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Logged out successfully",
	})
}
