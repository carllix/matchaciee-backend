package mocks

import (
	"github.com/carllix/matchaciee-backend/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) Create(token *models.RefreshToken) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) FindByToken(token string) (*models.RefreshToken, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	refreshToken, ok := args.Get(0).(*models.RefreshToken)
	if !ok {
		return nil, args.Error(1)
	}
	return refreshToken, args.Error(1)
}

func (m *MockRefreshTokenRepository) FindValidByToken(token string) (*models.RefreshToken, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	refreshToken, ok := args.Get(0).(*models.RefreshToken)
	if !ok {
		return nil, args.Error(1)
	}
	return refreshToken, args.Error(1)
}

func (m *MockRefreshTokenRepository) FindAllByUserID(userID uint) ([]*models.RefreshToken, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	tokens, ok := args.Get(0).([]*models.RefreshToken)
	if !ok {
		return nil, args.Error(1)
	}
	return tokens, args.Error(1)
}

func (m *MockRefreshTokenRepository) RevokeToken(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) RevokeAllUserTokens(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) DeleteExpiredTokens() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
