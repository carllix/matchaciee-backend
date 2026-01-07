package handlers

import (
	"errors"

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
