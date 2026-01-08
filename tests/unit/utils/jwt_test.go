package utils_test

import (
	"testing"
	"time"

	"github.com/carllix/matchaciee-backend/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJWTUtil(t *testing.T) {
	secretKey := "test-secret-key-at-least-32-characters-long"
	expiry := 1 * time.Hour
	refreshExpiry := 7 * 24 * time.Hour

	jwtUtil := utils.NewJWTUtil(secretKey, expiry, refreshExpiry)

	assert.NotNil(t, jwtUtil, "JWTUtil should not be nil")
}

func TestGenerateToken(t *testing.T) {
	secretKey := "test-secret-key-at-least-32-characters-long"
	expiry := 1 * time.Hour
	refreshExpiry := 7 * 24 * time.Hour
	jwtUtil := utils.NewJWTUtil(secretKey, expiry, refreshExpiry)

	t.Run("should generate valid token", func(t *testing.T) {
		userUUID := uuid.New()
		email := "test@example.com"
		role := "member"

		token, err := jwtUtil.GenerateToken(userUUID, email, role)

		require.NoError(t, err, "Should not return error")
		assert.NotEmpty(t, token, "Token should not be empty")

		// Verify token can be parsed
		claims, err := jwtUtil.ValidateToken(token)
		require.NoError(t, err, "Token should be valid")
		assert.Equal(t, userUUID, claims.UserUUID)
		assert.Equal(t, email, claims.Email)
		assert.Equal(t, role, claims.Role)
	})

	t.Run("should generate different tokens for different users", func(t *testing.T) {
		user1UUID := uuid.New()
		user2UUID := uuid.New()

		token1, err := jwtUtil.GenerateToken(user1UUID, "user1@example.com", "member")
		require.NoError(t, err)

		token2, err := jwtUtil.GenerateToken(user2UUID, "user2@example.com", "member")
		require.NoError(t, err)

		assert.NotEqual(t, token1, token2, "Tokens should be different")
	})

	t.Run("should include correct claims", func(t *testing.T) {
		userUUID := uuid.New()
		email := "test@example.com"
		role := "admin"

		token, err := jwtUtil.GenerateToken(userUUID, email, role)
		require.NoError(t, err)

		claims, err := jwtUtil.ValidateToken(token)
		require.NoError(t, err)

		assert.Equal(t, userUUID, claims.UserUUID)
		assert.Equal(t, email, claims.Email)
		assert.Equal(t, role, claims.Role)
		assert.NotNil(t, claims.ExpiresAt)
		assert.NotNil(t, claims.IssuedAt)
		assert.NotNil(t, claims.NotBefore)
	})
}

func TestValidateToken(t *testing.T) {
	secretKey := "test-secret-key-at-least-32-characters-long"
	expiry := 1 * time.Hour
	refreshExpiry := 7 * 24 * time.Hour
	jwtUtil := utils.NewJWTUtil(secretKey, expiry, refreshExpiry)

	t.Run("should validate valid token", func(t *testing.T) {
		userUUID := uuid.New()
		email := "test@example.com"
		role := "member"

		token, err := jwtUtil.GenerateToken(userUUID, email, role)
		require.NoError(t, err)

		claims, err := jwtUtil.ValidateToken(token)

		require.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, userUUID, claims.UserUUID)
		assert.Equal(t, email, claims.Email)
		assert.Equal(t, role, claims.Role)
	})

	t.Run("should reject invalid token", func(t *testing.T) {
		invalidToken := "invalid.token.here"

		claims, err := jwtUtil.ValidateToken(invalidToken)

		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.ErrorIs(t, err, utils.ErrInvalidToken)
	})

	t.Run("should reject empty token", func(t *testing.T) {
		claims, err := jwtUtil.ValidateToken("")

		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.ErrorIs(t, err, utils.ErrInvalidToken)
	})

	t.Run("should reject token with wrong secret", func(t *testing.T) {
		wrongSecretUtil := utils.NewJWTUtil("wrong-secret-key-different-from-original", expiry, refreshExpiry)
		userUUID := uuid.New()

		token, err := jwtUtil.GenerateToken(userUUID, "test@example.com", "member")
		require.NoError(t, err)

		claims, err := wrongSecretUtil.ValidateToken(token)

		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.ErrorIs(t, err, utils.ErrInvalidToken)
	})

	t.Run("should reject expired token", func(t *testing.T) {
		shortExpiryUtil := utils.NewJWTUtil(secretKey, 1*time.Millisecond, refreshExpiry)
		userUUID := uuid.New()

		token, err := shortExpiryUtil.GenerateToken(userUUID, "test@example.com", "member")
		require.NoError(t, err)

		// Wait for token to expire
		time.Sleep(10 * time.Millisecond)

		claims, err := shortExpiryUtil.ValidateToken(token)

		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.ErrorIs(t, err, utils.ErrExpiredToken)
	})

	t.Run("should reject malformed token", func(t *testing.T) {
		malformedTokens := []string{
			"not.a.token",
			"header.payload",
			"",
			".",
			"a.b.c.d",
		}

		for _, malformedToken := range malformedTokens {
			claims, err := jwtUtil.ValidateToken(malformedToken)

			assert.Error(t, err, "Should reject malformed token: %s", malformedToken)
			assert.Nil(t, claims)
			assert.ErrorIs(t, err, utils.ErrInvalidToken)
		}
	})

	t.Run("should reject token with invalid signing method", func(t *testing.T) {
		// Create token with RS256 instead of HS256
		userUUID := uuid.New()
		claims := utils.JWTClaims{
			UserUUID: userUUID,
			Email:    "test@example.com",
			Role:     "member",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
		tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
		require.NoError(t, err)

		validatedClaims, err := jwtUtil.ValidateToken(tokenString)

		assert.Error(t, err)
		assert.Nil(t, validatedClaims)
		assert.ErrorIs(t, err, utils.ErrInvalidToken)
	})
}

func TestGenerateRefreshToken(t *testing.T) {
	secretKey := "test-secret-key-at-least-32-characters-long"
	expiry := 1 * time.Hour
	refreshExpiry := 7 * 24 * time.Hour
	jwtUtil := utils.NewJWTUtil(secretKey, expiry, refreshExpiry)

	t.Run("should generate valid refresh token", func(t *testing.T) {
		userUUID := uuid.New()

		token, expiresAt, err := jwtUtil.GenerateRefreshToken(userUUID)

		require.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.False(t, expiresAt.IsZero())
		assert.True(t, expiresAt.After(time.Now()))
	})

	t.Run("should set correct expiration time", func(t *testing.T) {
		userUUID := uuid.New()
		beforeGeneration := time.Now()

		token, expiresAt, err := jwtUtil.GenerateRefreshToken(userUUID)

		require.NoError(t, err)
		assert.NotEmpty(t, token)

		expectedExpiry := beforeGeneration.Add(refreshExpiry)
		// Allow 1 second margin for test execution time
		assert.WithinDuration(t, expectedExpiry, expiresAt, 1*time.Second)
	})

	t.Run("should generate different tokens for same user", func(t *testing.T) {
		userUUID := uuid.New()

		token1, _, err := jwtUtil.GenerateRefreshToken(userUUID)
		require.NoError(t, err)

		// Small delay to ensure different issued time
		time.Sleep(1 * time.Second)

		token2, _, err := jwtUtil.GenerateRefreshToken(userUUID)
		require.NoError(t, err)

		assert.NotEqual(t, token1, token2, "Refresh tokens should be different even for same user")
	})

	t.Run("should generate different tokens for different users", func(t *testing.T) {
		user1UUID := uuid.New()
		user2UUID := uuid.New()

		token1, _, err := jwtUtil.GenerateRefreshToken(user1UUID)
		require.NoError(t, err)

		token2, _, err := jwtUtil.GenerateRefreshToken(user2UUID)
		require.NoError(t, err)

		assert.NotEqual(t, token1, token2)
	})
}

func TestValidateRefreshToken(t *testing.T) {
	secretKey := "test-secret-key-at-least-32-characters-long"
	expiry := 1 * time.Hour
	refreshExpiry := 7 * 24 * time.Hour
	jwtUtil := utils.NewJWTUtil(secretKey, expiry, refreshExpiry)

	t.Run("should validate valid refresh token", func(t *testing.T) {
		userUUID := uuid.New()

		token, _, err := jwtUtil.GenerateRefreshToken(userUUID)
		require.NoError(t, err)

		validatedUUID, err := jwtUtil.ValidateRefreshToken(token)

		require.NoError(t, err)
		assert.Equal(t, userUUID, validatedUUID)
	})

	t.Run("should reject invalid refresh token", func(t *testing.T) {
		invalidToken := "invalid.refresh.token"

		validatedUUID, err := jwtUtil.ValidateRefreshToken(invalidToken)

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, validatedUUID)
		assert.ErrorIs(t, err, utils.ErrInvalidToken)
	})

	t.Run("should reject empty refresh token", func(t *testing.T) {
		validatedUUID, err := jwtUtil.ValidateRefreshToken("")

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, validatedUUID)
		assert.ErrorIs(t, err, utils.ErrInvalidToken)
	})

	t.Run("should reject expired refresh token", func(t *testing.T) {
		shortRefreshUtil := utils.NewJWTUtil(secretKey, expiry, 1*time.Millisecond)
		userUUID := uuid.New()

		token, _, err := shortRefreshUtil.GenerateRefreshToken(userUUID)
		require.NoError(t, err)

		// Wait for token to expire
		time.Sleep(10 * time.Millisecond)

		validatedUUID, err := shortRefreshUtil.ValidateRefreshToken(token)

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, validatedUUID)
		assert.ErrorIs(t, err, utils.ErrExpiredToken)
	})

	t.Run("should reject refresh token with wrong secret", func(t *testing.T) {
		wrongSecretUtil := utils.NewJWTUtil("wrong-secret-key-different-from-original", expiry, refreshExpiry)
		userUUID := uuid.New()

		token, _, err := jwtUtil.GenerateRefreshToken(userUUID)
		require.NoError(t, err)

		validatedUUID, err := wrongSecretUtil.ValidateRefreshToken(token)

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, validatedUUID)
		assert.ErrorIs(t, err, utils.ErrInvalidToken)
	})

	t.Run("should reject refresh token with invalid UUID in subject", func(t *testing.T) {
		// Create token with invalid UUID in subject
		claims := jwt.RegisteredClaims{
			Subject:   "not-a-valid-uuid",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(secretKey))
		require.NoError(t, err)

		validatedUUID, err := jwtUtil.ValidateRefreshToken(tokenString)

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, validatedUUID)
		assert.ErrorIs(t, err, utils.ErrInvalidToken)
	})
}

func TestJWTTokenExpiry(t *testing.T) {
	secretKey := "test-secret-key-at-least-32-characters-long"

	t.Run("token should expire after configured duration", func(t *testing.T) {
		expiry := 2 * time.Second
		jwtUtil := utils.NewJWTUtil(secretKey, expiry, 1*time.Hour)
		userUUID := uuid.New()

		token, err := jwtUtil.GenerateToken(userUUID, "test@example.com", "member")
		require.NoError(t, err)

		// Token should be valid immediately
		claims, err := jwtUtil.ValidateToken(token)
		require.NoError(t, err)
		assert.NotNil(t, claims)

		// Wait for token to expire
		time.Sleep(3 * time.Second)

		// Token should now be expired
		claims, err = jwtUtil.ValidateToken(token)
		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.ErrorIs(t, err, utils.ErrExpiredToken)
	})
}
