package utils_test

import (
	"strings"
	"testing"

	"github.com/carllix/matchaciee-backend/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	t.Run("should hash password successfully", func(t *testing.T) {
		password := "SecurePassword123!"

		hashedPassword, err := utils.HashPassword(password)

		require.NoError(t, err, "Should not return error")
		assert.NotEmpty(t, hashedPassword, "Hashed password should not be empty")
		assert.NotEqual(t, password, hashedPassword, "Hashed password should be different from plain password")
	})

	t.Run("should generate different hashes for same password", func(t *testing.T) {
		password := "SamePassword123!"

		hash1, err := utils.HashPassword(password)
		require.NoError(t, err)

		hash2, err := utils.HashPassword(password)
		require.NoError(t, err)

		assert.NotEqual(t, hash1, hash2, "Different hashes should be generated due to random salt")
	})

	t.Run("should hash long password", func(t *testing.T) {
		password := strings.Repeat("a", 72)

		hashedPassword, err := utils.HashPassword(password)

		require.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
	})

	t.Run("should hash password with special characters", func(t *testing.T) {
		password := "P@ssw0rd!#$%^&*()_+-=[]{}|;':\",./<>?"

		hashedPassword, err := utils.HashPassword(password)

		require.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
	})

	t.Run("should hash password with unicode characters", func(t *testing.T) {
		password := "パスワード123!密码"

		hashedPassword, err := utils.HashPassword(password)

		require.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
	})

	t.Run("hashed password should be bcrypt format", func(t *testing.T) {
		password := "TestPassword123!"

		hashedPassword, err := utils.HashPassword(password)

		require.NoError(t, err)
		// Bcrypt hashes start with $2a$, $2b$, or $2y$ followed by cost
		assert.True(t, strings.HasPrefix(hashedPassword, "$2a$") ||
			strings.HasPrefix(hashedPassword, "$2b$") ||
			strings.HasPrefix(hashedPassword, "$2y$"),
			"Hash should be in bcrypt format")
	})

	t.Run("should use default cost", func(t *testing.T) {
		password := "TestPassword123!"

		hashedPassword, err := utils.HashPassword(password)

		require.NoError(t, err)

		cost, err := bcrypt.Cost([]byte(hashedPassword))
		require.NoError(t, err)
		assert.Equal(t, utils.DefaultCost, cost, "Should use DefaultCost (12)")
	})
}

func TestComparePassword(t *testing.T) {
	t.Run("should validate correct password", func(t *testing.T) {
		password := "CorrectPassword123!"
		hashedPassword, err := utils.HashPassword(password)
		require.NoError(t, err)

		err = utils.ComparePassword(hashedPassword, password)

		assert.NoError(t, err, "Should validate correct password")
	})

	t.Run("should reject incorrect password", func(t *testing.T) {
		password := "CorrectPassword123!"
		wrongPassword := "WrongPassword123!"
		hashedPassword, err := utils.HashPassword(password)
		require.NoError(t, err)

		err = utils.ComparePassword(hashedPassword, wrongPassword)

		assert.Error(t, err, "Should reject incorrect password")
		assert.ErrorIs(t, err, utils.ErrInvalidPassword)
	})

	t.Run("should be case sensitive", func(t *testing.T) {
		password := "Password123!"
		hashedPassword, err := utils.HashPassword(password)
		require.NoError(t, err)

		err = utils.ComparePassword(hashedPassword, "password123!")
		assert.Error(t, err)
		assert.ErrorIs(t, err, utils.ErrInvalidPassword)

		err = utils.ComparePassword(hashedPassword, "PASSWORD123!")
		assert.Error(t, err)
		assert.ErrorIs(t, err, utils.ErrInvalidPassword)
	})

	t.Run("should reject password with extra characters", func(t *testing.T) {
		password := "Password123!"
		hashedPassword, err := utils.HashPassword(password)
		require.NoError(t, err)

		err = utils.ComparePassword(hashedPassword, "Password123!extra")

		assert.Error(t, err)
		assert.ErrorIs(t, err, utils.ErrInvalidPassword)
	})

	t.Run("should reject password with missing characters", func(t *testing.T) {
		password := "Password123!"
		hashedPassword, err := utils.HashPassword(password)
		require.NoError(t, err)

		err = utils.ComparePassword(hashedPassword, "Password123")

		assert.Error(t, err)
		assert.ErrorIs(t, err, utils.ErrInvalidPassword)
	})

	t.Run("should reject invalid hash format", func(t *testing.T) {
		invalidHash := "not-a-valid-bcrypt-hash"
		password := "Password123!"

		err := utils.ComparePassword(invalidHash, password)

		assert.Error(t, err)
		assert.NotErrorIs(t, err, utils.ErrInvalidPassword)
	})

	t.Run("should validate password with special characters", func(t *testing.T) {
		password := "P@ssw0rd!#$%^&*()"
		hashedPassword, err := utils.HashPassword(password)
		require.NoError(t, err)

		err = utils.ComparePassword(hashedPassword, password)

		assert.NoError(t, err)
	})

	t.Run("should validate password with unicode characters", func(t *testing.T) {
		password := "パスワード123!密码"
		hashedPassword, err := utils.HashPassword(password)
		require.NoError(t, err)

		err = utils.ComparePassword(hashedPassword, password)

		assert.NoError(t, err)
	})

	t.Run("should validate long password", func(t *testing.T) {
		password := strings.Repeat("a", 72)
		hashedPassword, err := utils.HashPassword(password)
		require.NoError(t, err)

		err = utils.ComparePassword(hashedPassword, password)

		assert.NoError(t, err)
	})
}

func TestPasswordEdgeCases(t *testing.T) {
	t.Run("should handle password with whitespace", func(t *testing.T) {
		password := "  Password With Spaces  "

		hashedPassword, err := utils.HashPassword(password)
		require.NoError(t, err)

		// Exact match should work
		err = utils.ComparePassword(hashedPassword, password)
		assert.NoError(t, err)

		// Trimmed password should not work
		err = utils.ComparePassword(hashedPassword, strings.TrimSpace(password))
		assert.Error(t, err)
		assert.ErrorIs(t, err, utils.ErrInvalidPassword)
	})

	t.Run("should handle password with newlines and tabs", func(t *testing.T) {
		password := "Pass\nword\t123!"

		hashedPassword, err := utils.HashPassword(password)
		require.NoError(t, err)

		err = utils.ComparePassword(hashedPassword, password)
		assert.NoError(t, err)
	})
}

func BenchmarkHashPassword(b *testing.B) {
	password := "BenchmarkPassword123!"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = utils.HashPassword(password)
	}
}

func BenchmarkComparePassword(b *testing.B) {
	password := "BenchmarkPassword123!"
	hashedPassword, _ := utils.HashPassword(password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = utils.ComparePassword(hashedPassword, password)
	}
}
