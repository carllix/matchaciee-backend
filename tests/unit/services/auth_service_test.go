package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/carllix/matchaciee-backend/internal/models"
	"github.com/carllix/matchaciee-backend/internal/repositories"
	"github.com/carllix/matchaciee-backend/internal/services"
	"github.com/carllix/matchaciee-backend/internal/utils"
	"github.com/carllix/matchaciee-backend/tests/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupAuthServiceTest() (
	*mocks.MockUserRepository,
	*mocks.MockRefreshTokenRepository,
	*utils.JWTUtil,
	services.AuthService,
) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockRefreshTokenRepo := new(mocks.MockRefreshTokenRepository)
	jwtUtil := utils.NewJWTUtil(
		"test-secret-key-at-least-32-characters-long",
		1*time.Hour,
		7*24*time.Hour,
	)

	authService := services.NewAuthService(mockUserRepo, mockRefreshTokenRepo, jwtUtil)

	return mockUserRepo, mockRefreshTokenRepo, jwtUtil, authService
}

func TestRegister(t *testing.T) {
	t.Run("should register user successfully", func(t *testing.T) {
		mockUserRepo, mockRefreshTokenRepo, _, authService := setupAuthServiceTest()

		req := services.RegisterRequest{
			Email:    "test@example.com",
			Password: "SecurePassword123!",
			FullName: "Test User",
			Phone:    "+6281234567890",
		}

		mockUserRepo.On("Create", mock.MatchedBy(func(user *models.User) bool {
			return user.Email == req.Email &&
				user.FullName == req.FullName &&
				user.Role == models.RoleMember &&
				user.IsActive == true
		})).Run(func(args mock.Arguments) {
			user := args.Get(0).(*models.User)
			user.ID = 1
			user.UUID = uuid.New()
		}).Return(nil)

		mockRefreshTokenRepo.On("Create", mock.AnythingOfType("*models.RefreshToken")).Return(nil)

		resp, err := authService.Register(req)

		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Token)
		assert.NotEmpty(t, resp.RefreshToken)
		assert.Equal(t, req.Email, resp.User.Email)
		assert.Equal(t, req.FullName, resp.User.FullName)
		assert.Equal(t, models.RoleMember, resp.User.Role)

		mockUserRepo.AssertExpectations(t)
		mockRefreshTokenRepo.AssertExpectations(t)
	})

	t.Run("should return error when email already exists", func(t *testing.T) {
		mockUserRepo, mockRefreshTokenRepo, _, authService := setupAuthServiceTest()

		req := services.RegisterRequest{
			Email:    "existing@example.com",
			Password: "SecurePassword123!",
			FullName: "Test User",
		}

		mockUserRepo.On("Create", mock.AnythingOfType("*models.User")).
			Return(repositories.ErrEmailAlreadyExists)

		resp, err := authService.Register(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.ErrorIs(t, err, repositories.ErrEmailAlreadyExists)

		mockUserRepo.AssertExpectations(t)
		mockRefreshTokenRepo.AssertNotCalled(t, "Create")
	})

	t.Run("should handle database error during user creation", func(t *testing.T) {
		mockUserRepo, mockRefreshTokenRepo, _, authService := setupAuthServiceTest()

		req := services.RegisterRequest{
			Email:    "test@example.com",
			Password: "SecurePassword123!",
			FullName: "Test User",
		}

		dbError := errors.New("database connection failed")
		mockUserRepo.On("Create", mock.AnythingOfType("*models.User")).Return(dbError)

		resp, err := authService.Register(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, dbError, err)

		mockUserRepo.AssertExpectations(t)
		mockRefreshTokenRepo.AssertNotCalled(t, "Create")
	})

	t.Run("should handle error when creating refresh token", func(t *testing.T) {
		mockUserRepo, mockRefreshTokenRepo, _, authService := setupAuthServiceTest()

		req := services.RegisterRequest{
			Email:    "test@example.com",
			Password: "SecurePassword123!",
			FullName: "Test User",
		}

		mockUserRepo.On("Create", mock.AnythingOfType("*models.User")).Run(func(args mock.Arguments) {
			user := args.Get(0).(*models.User)
			user.ID = 1
			user.UUID = uuid.New()
		}).Return(nil)

		tokenError := errors.New("failed to create refresh token")
		mockRefreshTokenRepo.On("Create", mock.AnythingOfType("*models.RefreshToken")).Return(tokenError)

		resp, err := authService.Register(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, tokenError, err)

		mockUserRepo.AssertExpectations(t)
		mockRefreshTokenRepo.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	t.Run("should login successfully with valid credentials", func(t *testing.T) {
		mockUserRepo, mockRefreshTokenRepo, _, authService := setupAuthServiceTest()

		password := "SecurePassword123!"
		hashedPassword, _ := utils.HashPassword(password)

		existingUser := &models.User{
			ID:       1,
			UUID:     uuid.New(),
			Email:    "test@example.com",
			Password: hashedPassword,
			FullName: "Test User",
			Role:     models.RoleMember,
			IsActive: true,
		}

		req := services.LoginRequest{
			Email:    "test@example.com",
			Password: password,
		}

		mockUserRepo.On("FindByEmail", req.Email).Return(existingUser, nil)
		mockRefreshTokenRepo.On("Create", mock.AnythingOfType("*models.RefreshToken")).Return(nil)

		resp, err := authService.Login(req)

		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Token)
		assert.NotEmpty(t, resp.RefreshToken)
		assert.Equal(t, existingUser.Email, resp.User.Email)
		assert.Equal(t, existingUser.FullName, resp.User.FullName)

		mockUserRepo.AssertExpectations(t)
		mockRefreshTokenRepo.AssertExpectations(t)
	})

	t.Run("should return error with invalid email", func(t *testing.T) {
		mockUserRepo, mockRefreshTokenRepo, _, authService := setupAuthServiceTest()

		req := services.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "SecurePassword123!",
		}

		mockUserRepo.On("FindByEmail", req.Email).Return(nil, repositories.ErrUserNotFound)

		resp, err := authService.Login(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.ErrorIs(t, err, services.ErrInvalidCredentials)

		mockUserRepo.AssertExpectations(t)
		mockRefreshTokenRepo.AssertNotCalled(t, "Create")
	})

	t.Run("should return error with invalid password", func(t *testing.T) {
		mockUserRepo, mockRefreshTokenRepo, _, authService := setupAuthServiceTest()

		correctPassword := "CorrectPassword123!"
		hashedPassword, _ := utils.HashPassword(correctPassword)

		existingUser := &models.User{
			ID:       1,
			UUID:     uuid.New(),
			Email:    "test@example.com",
			Password: hashedPassword,
			FullName: "Test User",
			Role:     models.RoleMember,
			IsActive: true,
		}

		req := services.LoginRequest{
			Email:    "test@example.com",
			Password: "WrongPassword123!",
		}

		mockUserRepo.On("FindByEmail", req.Email).Return(existingUser, nil)

		resp, err := authService.Login(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.ErrorIs(t, err, services.ErrInvalidCredentials)

		mockUserRepo.AssertExpectations(t)
		mockRefreshTokenRepo.AssertNotCalled(t, "Create")
	})

	t.Run("should return error when user is inactive", func(t *testing.T) {
		mockUserRepo, mockRefreshTokenRepo, _, authService := setupAuthServiceTest()

		password := "SecurePassword123!"
		hashedPassword, _ := utils.HashPassword(password)

		inactiveUser := &models.User{
			ID:       1,
			UUID:     uuid.New(),
			Email:    "test@example.com",
			Password: hashedPassword,
			FullName: "Test User",
			Role:     models.RoleMember,
			IsActive: false, // User is inactive
		}

		req := services.LoginRequest{
			Email:    "test@example.com",
			Password: password,
		}

		mockUserRepo.On("FindByEmail", req.Email).Return(inactiveUser, nil)

		resp, err := authService.Login(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.ErrorIs(t, err, services.ErrUserInactive)

		mockUserRepo.AssertExpectations(t)
		mockRefreshTokenRepo.AssertNotCalled(t, "Create")
	})

	t.Run("should handle database error during login", func(t *testing.T) {
		mockUserRepo, mockRefreshTokenRepo, _, authService := setupAuthServiceTest()

		req := services.LoginRequest{
			Email:    "test@example.com",
			Password: "SecurePassword123!",
		}

		dbError := errors.New("database error")
		mockUserRepo.On("FindByEmail", req.Email).Return(nil, dbError)

		resp, err := authService.Login(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, dbError, err)

		mockUserRepo.AssertExpectations(t)
		mockRefreshTokenRepo.AssertNotCalled(t, "Create")
	})
}

func TestRefreshToken(t *testing.T) {
	t.Run("should refresh token successfully", func(t *testing.T) {
		mockUserRepo, mockRefreshTokenRepo, jwtUtil, authService := setupAuthServiceTest()

		userUUID := uuid.New()
		oldRefreshToken, expiresAt, _ := jwtUtil.GenerateRefreshToken(userUUID)

		existingUser := &models.User{
			ID:       1,
			UUID:     userUUID,
			Email:    "test@example.com",
			FullName: "Test User",
			Role:     models.RoleMember,
			IsActive: true,
		}

		refreshTokenModel := &models.RefreshToken{
			ID:        1,
			UserID:    1,
			Token:     oldRefreshToken,
			ExpiresAt: expiresAt,
			RevokedAt: nil,
		}

		req := services.RefreshTokenRequest{
			RefreshToken: oldRefreshToken,
		}

		mockRefreshTokenRepo.On("FindValidByToken", oldRefreshToken).Return(refreshTokenModel, nil)
		mockUserRepo.On("FindByID", uint(1)).Return(existingUser, nil)
		mockRefreshTokenRepo.On("RevokeToken", oldRefreshToken).Return(nil)
		mockRefreshTokenRepo.On("Create", mock.AnythingOfType("*models.RefreshToken")).Return(nil)

		// Small delay to ensure different issued time
		time.Sleep(1 * time.Second)

		resp, err := authService.RefreshToken(req)

		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Token)
		assert.NotEmpty(t, resp.RefreshToken)
		assert.NotEqual(t, oldRefreshToken, resp.RefreshToken) // Should be new token

		mockUserRepo.AssertExpectations(t)
		mockRefreshTokenRepo.AssertExpectations(t)
	})

	t.Run("should return error with invalid refresh token", func(t *testing.T) {
		mockUserRepo, mockRefreshTokenRepo, _, authService := setupAuthServiceTest()

		req := services.RefreshTokenRequest{
			RefreshToken: "invalid-token",
		}

		resp, err := authService.RefreshToken(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.ErrorIs(t, err, services.ErrInvalidRefreshToken)

		mockUserRepo.AssertNotCalled(t, "FindByID")
		mockRefreshTokenRepo.AssertNotCalled(t, "FindValidByToken")
	})

	t.Run("should return error when refresh token not found in database", func(t *testing.T) {
		mockUserRepo, mockRefreshTokenRepo, jwtUtil, authService := setupAuthServiceTest()

		userUUID := uuid.New()
		validToken, _, _ := jwtUtil.GenerateRefreshToken(userUUID)

		req := services.RefreshTokenRequest{
			RefreshToken: validToken,
		}

		dbError := errors.New("token not found")
		mockRefreshTokenRepo.On("FindValidByToken", validToken).Return(nil, dbError)

		resp, err := authService.RefreshToken(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.ErrorIs(t, err, services.ErrInvalidRefreshToken)

		mockRefreshTokenRepo.AssertExpectations(t)
		mockUserRepo.AssertNotCalled(t, "FindByID")
	})

	t.Run("should return error when refresh token is expired", func(t *testing.T) {
		mockUserRepo, mockRefreshTokenRepo, jwtUtil, authService := setupAuthServiceTest()

		userUUID := uuid.New()
		oldToken, _, _ := jwtUtil.GenerateRefreshToken(userUUID)

		// Expired token
		expiredTokenModel := &models.RefreshToken{
			ID:        1,
			UserID:    1,
			Token:     oldToken,
			ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired
			RevokedAt: nil,
		}

		req := services.RefreshTokenRequest{
			RefreshToken: oldToken,
		}

		mockRefreshTokenRepo.On("FindValidByToken", oldToken).Return(expiredTokenModel, nil)

		resp, err := authService.RefreshToken(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.ErrorIs(t, err, services.ErrInvalidRefreshToken)

		mockRefreshTokenRepo.AssertExpectations(t)
		mockUserRepo.AssertNotCalled(t, "FindByID")
	})

	t.Run("should return error when user is inactive", func(t *testing.T) {
		mockUserRepo, mockRefreshTokenRepo, jwtUtil, authService := setupAuthServiceTest()

		userUUID := uuid.New()
		validToken, expiresAt, _ := jwtUtil.GenerateRefreshToken(userUUID)

		inactiveUser := &models.User{
			ID:       1,
			UUID:     userUUID,
			Email:    "test@example.com",
			FullName: "Test User",
			Role:     models.RoleMember,
			IsActive: false, // Inactive
		}

		refreshTokenModel := &models.RefreshToken{
			ID:        1,
			UserID:    1,
			Token:     validToken,
			ExpiresAt: expiresAt,
			RevokedAt: nil,
		}

		req := services.RefreshTokenRequest{
			RefreshToken: validToken,
		}

		mockRefreshTokenRepo.On("FindValidByToken", validToken).Return(refreshTokenModel, nil)
		mockUserRepo.On("FindByID", uint(1)).Return(inactiveUser, nil)

		resp, err := authService.RefreshToken(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.ErrorIs(t, err, services.ErrUserInactive)

		mockUserRepo.AssertExpectations(t)
		mockRefreshTokenRepo.AssertExpectations(t)
		mockRefreshTokenRepo.AssertNotCalled(t, "RevokeToken")
	})
}

func TestLogout(t *testing.T) {
	t.Run("should logout successfully", func(t *testing.T) {
		mockUserRepo, mockRefreshTokenRepo, jwtUtil, authService := setupAuthServiceTest()

		userUUID := uuid.New()
		refreshToken, _, _ := jwtUtil.GenerateRefreshToken(userUUID)

		mockRefreshTokenRepo.On("RevokeToken", refreshToken).Return(nil)

		err := authService.Logout(refreshToken)

		assert.NoError(t, err)

		mockRefreshTokenRepo.AssertExpectations(t)
		mockUserRepo.AssertNotCalled(t, "FindByID")
	})

	t.Run("should return error with invalid token", func(t *testing.T) {
		_, _, _, authService := setupAuthServiceTest()

		invalidToken := "invalid-token"

		err := authService.Logout(invalidToken)

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInvalidRefreshToken)
	})

	t.Run("should handle error during token revocation", func(t *testing.T) {
		mockUserRepo, mockRefreshTokenRepo, jwtUtil, authService := setupAuthServiceTest()

		userUUID := uuid.New()
		refreshToken, _, _ := jwtUtil.GenerateRefreshToken(userUUID)

		dbError := errors.New("database error")
		mockRefreshTokenRepo.On("RevokeToken", refreshToken).Return(dbError)

		err := authService.Logout(refreshToken)

		assert.Error(t, err)
		assert.Equal(t, dbError, err)

		mockRefreshTokenRepo.AssertExpectations(t)
		mockUserRepo.AssertNotCalled(t, "FindByID")
	})
}

func TestGetUserByUUID(t *testing.T) {
	t.Run("should get user successfully", func(t *testing.T) {
		mockUserRepo, mockRefreshTokenRepo, _, authService := setupAuthServiceTest()

		userUUID := uuid.New()
		existingUser := &models.User{
			ID:       1,
			UUID:     userUUID,
			Email:    "test@example.com",
			FullName: "Test User",
			Role:     models.RoleMember,
			IsActive: true,
		}

		mockUserRepo.On("FindByUUID", userUUID).Return(existingUser, nil)

		resp, err := authService.GetUserByUUID(userUUID)

		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, userUUID, resp.ID)
		assert.Equal(t, existingUser.Email, resp.Email)
		assert.Equal(t, existingUser.FullName, resp.FullName)
		assert.Equal(t, existingUser.Role, resp.Role)

		mockUserRepo.AssertExpectations(t)
		mockRefreshTokenRepo.AssertNotCalled(t, "Create")
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		mockUserRepo, mockRefreshTokenRepo, _, authService := setupAuthServiceTest()

		userUUID := uuid.New()

		mockUserRepo.On("FindByUUID", userUUID).Return(nil, repositories.ErrUserNotFound)

		resp, err := authService.GetUserByUUID(userUUID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.ErrorIs(t, err, repositories.ErrUserNotFound)

		mockUserRepo.AssertExpectations(t)
		mockRefreshTokenRepo.AssertNotCalled(t, "Create")
	})

	t.Run("should handle database error", func(t *testing.T) {
		mockUserRepo, mockRefreshTokenRepo, _, authService := setupAuthServiceTest()

		userUUID := uuid.New()
		dbError := errors.New("database error")

		mockUserRepo.On("FindByUUID", userUUID).Return(nil, dbError)

		resp, err := authService.GetUserByUUID(userUUID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, dbError, err)

		mockUserRepo.AssertExpectations(t)
		mockRefreshTokenRepo.AssertNotCalled(t, "Create")
	})
}
