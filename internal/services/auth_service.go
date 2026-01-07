package services

import (
	"errors"

	"github.com/carllix/matchaciee-backend/internal/models"
	"github.com/carllix/matchaciee-backend/internal/repositories"
	"github.com/carllix/matchaciee-backend/internal/utils"
	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrUserInactive        = errors.New("user account is inactive")
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required,min=2"`
	Phone    string `json:"phone,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token        string       `json:"token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type UserResponse struct {
	Phone    *string         `json:"phone,omitempty"`
	Email    string          `json:"email"`
	FullName string          `json:"full_name"`
	Role     models.UserRole `json:"role"`
	ID       uuid.UUID       `json:"id"`
}

type AuthService interface {
	Register(req RegisterRequest) (*AuthResponse, error)
	Login(req LoginRequest) (*AuthResponse, error)
	RefreshToken(req RefreshTokenRequest) (*AuthResponse, error)
	Logout(refreshToken string) error
	GetUserByUUID(uuid uuid.UUID) (*UserResponse, error)
}

type authService struct {
	userRepo         repositories.UserRepository
	refreshTokenRepo repositories.RefreshTokenRepository
	jwtUtil          *utils.JWTUtil
}

func NewAuthService(
	userRepo repositories.UserRepository,
	refreshTokenRepo repositories.RefreshTokenRepository,
	jwtUtil *utils.JWTUtil,
) AuthService {
	return &authService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtUtil:          jwtUtil,
	}
}

func (s *authService) Register(req RegisterRequest) (*AuthResponse, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:    req.Email,
		Password: hashedPassword,
		FullName: req.FullName,
		Role:     models.RoleMember,
		IsActive: true,
	}

	if req.Phone != "" {
		user.Phone = &req.Phone
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	token, err := s.jwtUtil.GenerateToken(user.UUID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	refreshToken, expiresAt, err := s.jwtUtil.GenerateRefreshToken(user.UUID)
	if err != nil {
		return nil, err
	}

	refreshTokenModel := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: expiresAt,
	}
	err = s.refreshTokenRepo.Create(refreshTokenModel)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         s.toUserResponse(user),
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) Login(req LoginRequest) (*AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	err = utils.ComparePassword(user.Password, req.Password)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.jwtUtil.GenerateToken(user.UUID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	refreshToken, expiresAt, err := s.jwtUtil.GenerateRefreshToken(user.UUID)
	if err != nil {
		return nil, err
	}

	refreshTokenModel := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: expiresAt,
	}
	err = s.refreshTokenRepo.Create(refreshTokenModel)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         s.toUserResponse(user),
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) RefreshToken(req RefreshTokenRequest) (*AuthResponse, error) {
	userUUID, err := s.jwtUtil.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	refreshToken, err := s.refreshTokenRepo.FindValidByToken(req.RefreshToken)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	if !refreshToken.IsValid() {
		return nil, ErrInvalidRefreshToken
	}

	user, err := s.userRepo.FindByID(refreshToken.UserID)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	accessToken, err := s.jwtUtil.GenerateToken(user.UUID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	// Generate new refresh token (token rotation)
	newRefreshToken, expiresAt, err := s.jwtUtil.GenerateRefreshToken(userUUID)
	if err != nil {
		return nil, err
	}

	// Revoke old refresh token
	err = s.refreshTokenRepo.RevokeToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Store new refresh token
	newRefreshTokenModel := &models.RefreshToken{
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiresAt: expiresAt,
	}
	err = s.refreshTokenRepo.Create(newRefreshTokenModel)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         s.toUserResponse(user),
		Token:        accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *authService) Logout(refreshToken string) error {
	_, err := s.jwtUtil.ValidateRefreshToken(refreshToken)
	if err != nil {
		return ErrInvalidRefreshToken
	}

	err = s.refreshTokenRepo.RevokeToken(refreshToken)
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) GetUserByUUID(uuid uuid.UUID) (*UserResponse, error) {
	user, err := s.userRepo.FindByUUID(uuid)
	if err != nil {
		return nil, err
	}

	userResp := s.toUserResponse(user)
	return &userResp, nil
}

func (s *authService) toUserResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:       user.UUID,
		Email:    user.Email,
		FullName: user.FullName,
		Phone:    user.Phone,
		Role:     user.Role,
	}
}
