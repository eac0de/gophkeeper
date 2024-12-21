package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eac0de/gophkeeper/auth/pkg/outmiddlewares"
	"github.com/eac0de/gophkeeper/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockIUserDataService is a mock implementation of IUserDataService for testing purposes.
type MockIUserDataService struct {
	mock.Mock
}

func (m *MockIUserDataService) InsertUserTextData(ctx context.Context, userID uuid.UUID, name string, text string, metadata map[string]interface{}) (*models.UserTextData, error) {
	args := m.Called(ctx, userID, name, text, metadata)
	return args.Get(0).(*models.UserTextData), args.Error(1)
}

func (m *MockIUserDataService) InsertUserFileData(ctx context.Context, userID uuid.UUID, name string, pathToFile string) (*models.UserFileData, error) {
	args := m.Called(ctx, userID, name, pathToFile)
	return args.Get(0).(*models.UserFileData), args.Error(1)
}

func (m *MockIUserDataService) InsertUserAuthInfo(ctx context.Context, userID uuid.UUID, name, login, password string, metadata map[string]interface{}) (*models.UserAuthInfo, error) {
	args := m.Called(ctx, userID, name, login, password, metadata)
	return args.Get(0).(*models.UserAuthInfo), args.Error(1)
}

func (m *MockIUserDataService) InsertUserBankCard(ctx context.Context, userID uuid.UUID, name, number, cardHolder, expireDate, csc string, metadata map[string]interface{}) (*models.UserBankCard, error) {
	args := m.Called(ctx, userID, name, number, cardHolder, expireDate, csc, metadata)
	return args.Get(0).(*models.UserBankCard), args.Error(1)
}

func (m *MockIUserDataService) UpdateUserTextData(ctx context.Context, userID uuid.UUID, ID uuid.UUID, name, text string, metadata map[string]interface{}) (*models.UserTextData, error) {
	args := m.Called(ctx, userID, ID, name, text, metadata)
	return args.Get(0).(*models.UserTextData), args.Error(1)
}

func (m *MockIUserDataService) UpdateUserFileData(ctx context.Context, userID uuid.UUID, ID uuid.UUID, name string, metadata map[string]interface{}) (*models.UserFileData, error) {
	args := m.Called(ctx, userID, ID, name, metadata)
	return args.Get(0).(*models.UserFileData), args.Error(1)
}

func (m *MockIUserDataService) UpdateUserAuthInfo(ctx context.Context, userID uuid.UUID, ID uuid.UUID, name, login, password string, metadata map[string]interface{}) (*models.UserAuthInfo, error) {
	args := m.Called(ctx, userID, ID, name, login, password, metadata)
	return args.Get(0).(*models.UserAuthInfo), args.Error(1)
}

func (m *MockIUserDataService) UpdateUserBankCard(ctx context.Context, userID uuid.UUID, ID uuid.UUID, name, number, cardHolder, expireDate, csc string, metadata map[string]interface{}) (*models.UserBankCard, error) {
	args := m.Called(ctx, userID, ID, name, number, cardHolder, expireDate, csc, metadata)
	return args.Get(0).(*models.UserBankCard), args.Error(1)
}

func (m *MockIUserDataService) GetUserTextData(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) (*models.UserTextData, error) {
	args := m.Called(ctx, dataID, userID)
	return args.Get(0).(*models.UserTextData), args.Error(1)
}

func (m *MockIUserDataService) GetUserFileData(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) (*models.UserFileData, error) {
	args := m.Called(ctx, dataID, userID)
	return args.Get(0).(*models.UserFileData), args.Error(1)
}

func (m *MockIUserDataService) GetUserAuthInfo(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) (*models.UserAuthInfo, error) {
	args := m.Called(ctx, dataID, userID)
	return args.Get(0).(*models.UserAuthInfo), args.Error(1)
}

func (m *MockIUserDataService) GetUserBankCard(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) (*models.UserBankCard, error) {
	args := m.Called(ctx, dataID, userID)
	return args.Get(0).(*models.UserBankCard), args.Error(1)
}

func (m *MockIUserDataService) DeleteUserTextData(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) error {
	args := m.Called(ctx, dataID, userID)
	return args.Error(0)
}

func (m *MockIUserDataService) DeleteUserFileData(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) error {
	args := m.Called(ctx, dataID, userID)
	return args.Error(0)
}

func (m *MockIUserDataService) DeleteUserAuthInfo(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) error {
	args := m.Called(ctx, dataID, userID)
	return args.Error(0)
}

func (m *MockIUserDataService) DeleteUserBankCard(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) error {
	args := m.Called(ctx, dataID, userID)
	return args.Error(0)
}

func NewTestAuthMiddleware(userID uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(gin.AuthUserKey, userID)
		c.Next()
	}
}

func TestInsertUserAuthInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockIUserDataService)
	handlers := NewUserDataHandlers(mockService)

	router := gin.Default()
	rootGroup := router.Group("api/gophkeeper/")
	authenticatedGroup := rootGroup.Group("/", outmiddlewares.NewAuthMiddlewareForTest())
	authenticatedGroup.GET("/user_auth_info/:id/", handlers.GetUserAuthInfo)
	authenticatedGroup.DELETE("/user_auth_info/:id/", handlers.DeleteUserAuthInfo)
	authenticatedGroup.PUT("/user_auth_info/:id/", handlers.UpdateUserAuthInfo)
	authenticatedGroup.POST("/user_auth_info/", handlers.InsertUserAuthInfo)

	t.Run("Success", func(t *testing.T) {
		requestBody, _ := json.Marshal(gin.H{
			"name":     "testName",
			"login":    "testLogin",
			"password": "testPassword",
			"metadata": gin.H{},
		})

		req, _ := http.NewRequest(http.MethodPost, "/authinfo", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		mockService.On("InsertUserAuthInfo", mock.Anything, mock.Anything, "testName", "testLogin", "testPassword", mock.Anything).Return(&models.UserAuthInfo{}, nil)

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Missing Fields", func(t *testing.T) {
		requestBody, _ := json.Marshal(gin.H{
			"name":  "testName",
			"login": "testLogin",
			// "password" field is missing
			"metadata": gin.H{},
		})

		req, _ := http.NewRequest(http.MethodPost, "/authinfo", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
