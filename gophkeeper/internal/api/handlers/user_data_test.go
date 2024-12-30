package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/eac0de/gophkeeper/auth/pkg/outmiddlewares"
	"github.com/eac0de/gophkeeper/internal/models"
	"github.com/eac0de/gophkeeper/shared/pkg/httperror"
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
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserTextData), args.Error(1)
}

func (m *MockIUserDataService) InsertUserFileData(ctx context.Context, userID uuid.UUID, name string, pathToFile string) (*models.UserFileData, error) {
	args := m.Called(ctx, userID, name, pathToFile)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserFileData), args.Error(1)
}

func (m *MockIUserDataService) InsertUserAuthInfo(ctx context.Context, userID uuid.UUID, name, login, password string, metadata map[string]interface{}) (*models.UserAuthInfo, error) {
	args := m.Called(ctx, userID, name, login, password, metadata)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserAuthInfo), args.Error(1)
}

func (m *MockIUserDataService) InsertUserBankCard(ctx context.Context, userID uuid.UUID, name, number, cardHolder, expireDate, csc string, metadata map[string]interface{}) (*models.UserBankCard, error) {
	args := m.Called(ctx, userID, name, number, cardHolder, expireDate, csc, metadata)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserBankCard), args.Error(1)
}

func (m *MockIUserDataService) UpdateUserTextData(ctx context.Context, userID uuid.UUID, ID uuid.UUID, name, text string, metadata map[string]interface{}) (*models.UserTextData, error) {
	args := m.Called(ctx, userID, ID, name, text, metadata)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserTextData), args.Error(1)
}

func (m *MockIUserDataService) UpdateUserFileData(ctx context.Context, userID uuid.UUID, ID uuid.UUID, name string, metadata map[string]interface{}) (*models.UserFileData, error) {
	args := m.Called(ctx, userID, ID, name, metadata)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserFileData), args.Error(1)
}

func (m *MockIUserDataService) UpdateUserAuthInfo(ctx context.Context, userID uuid.UUID, ID uuid.UUID, name, login, password string, metadata map[string]interface{}) (*models.UserAuthInfo, error) {
	args := m.Called(ctx, userID, ID, name, login, password, metadata)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserAuthInfo), args.Error(1)
}

func (m *MockIUserDataService) UpdateUserBankCard(ctx context.Context, userID uuid.UUID, ID uuid.UUID, name, number, cardHolder, expireDate, csc string, metadata map[string]interface{}) (*models.UserBankCard, error) {
	args := m.Called(ctx, userID, ID, name, number, cardHolder, expireDate, csc, metadata)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserBankCard), args.Error(1)
}

func (m *MockIUserDataService) GetUserTextData(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) (*models.UserTextData, error) {
	args := m.Called(ctx, dataID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserTextData), args.Error(1)
}

func (m *MockIUserDataService) GetUserFileData(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) (*models.UserFileData, error) {
	args := m.Called(ctx, dataID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserFileData), args.Error(1)
}

func (m *MockIUserDataService) GetUserAuthInfo(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) (*models.UserAuthInfo, error) {
	args := m.Called(ctx, dataID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserAuthInfo), args.Error(1)
}

func (m *MockIUserDataService) GetUserBankCard(ctx context.Context, dataID uuid.UUID, userID uuid.UUID) (*models.UserBankCard, error) {
	args := m.Called(ctx, dataID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
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

func TestInsertUserAuthInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockIUserDataService)
	handlers := NewUserDataHandlers(mockService)

	router := gin.Default()
	userID := uuid.New()
	authenticatedGroup := router.Group("/", outmiddlewares.NewAuthMiddlewareForTest(userID))
	authenticatedGroup.POST("/user_auth_info/", handlers.InsertUserAuthInfo)

	t.Run("Success", func(t *testing.T) {
		requestBody, _ := json.Marshal(gin.H{
			"name":     "testName",
			"login":    "testLogin",
			"password": "testPassword",
			"metadata": gin.H{},
		})

		req, _ := http.NewRequest(http.MethodPost, "/user_auth_info/", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		mockService.On("InsertUserAuthInfo", mock.Anything, mock.Anything, "testName", "testLogin", "testPassword", mock.Anything).Return(&models.UserAuthInfo{}, nil).Once()

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

		req, _ := http.NewRequest(http.MethodPost, "/user_auth_info/", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		invalidJSON := `{"name":"testName", "login": "testLogin", "password"` // неправильный JSON
		req, _ := http.NewRequest(http.MethodPost, "/user_auth_info/", bytes.NewBufferString(invalidJSON))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "unexpected EOF")
	})

	t.Run("Service Error", func(t *testing.T) {
		requestBody, _ := json.Marshal(gin.H{
			"name":     "testName",
			"login":    "testLogin",
			"password": "testPassword",
			"metadata": gin.H{},
		})

		req, _ := http.NewRequest(http.MethodPost, "/user_auth_info/", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		simulatedError := errors.New("simulated service error")
		mockService.On("InsertUserAuthInfo", mock.Anything, mock.Anything, "testName", "testLogin", "testPassword", mock.Anything).Return(nil, simulatedError).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "simulated service error")
		mockService.AssertExpectations(t)
	})

	t.Run("Nil Fields", func(t *testing.T) {
		requestBody, _ := json.Marshal(gin.H{
			"name":     nil,
			"login":    "testLogin",
			"password": "testPassword",
			"metadata": gin.H{},
		})

		req, _ := http.NewRequest(http.MethodPost, "/user_auth_info/", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "name,login and password are required")
	})
}
func TestDeleteUserAuthInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockIUserDataService)
	handlers := NewUserDataHandlers(mockService)

	router := gin.Default()
	userID := uuid.New()
	authenticatedGroup := router.Group("/", outmiddlewares.NewAuthMiddlewareForTest(userID))
	authenticatedGroup.DELETE("/user_auth_info/:id/", handlers.DeleteUserAuthInfo)

	t.Run("Success", func(t *testing.T) {
		dataID := uuid.New()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/user_auth_info/%s/", dataID.String()), nil)

		rec := httptest.NewRecorder()

		mockService.On("DeleteUserAuthInfo", mock.Anything, dataID, userID).Return(nil).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		invalidID := "invalid-uuid"
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/user_auth_info/%s/", invalidID), nil)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid data id")
	})

	t.Run("Service Error", func(t *testing.T) {
		dataID := uuid.New()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/user_auth_info/%s/", dataID.String()), nil)

		rec := httptest.NewRecorder()

		simulatedError := errors.New("simulated service error")
		mockService.On("DeleteUserAuthInfo", mock.Anything, dataID, userID).Return(simulatedError).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "simulated service error")
		mockService.AssertExpectations(t)
	})
}

func TestUpdateUserAuthInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockIUserDataService)
	handlers := NewUserDataHandlers(mockService)

	router := gin.Default()
	userID := uuid.New()
	authenticatedGroup := router.Group("/", outmiddlewares.NewAuthMiddlewareForTest(userID))
	authenticatedGroup.PUT("/user_auth_info/:id/", handlers.UpdateUserAuthInfo)

	t.Run("Success", func(t *testing.T) {
		dataID := uuid.New()
		requestBody, _ := json.Marshal(gin.H{
			"name":     "testName",
			"login":    "testLogin",
			"password": "testPassword",
			"metadata": gin.H{},
		})

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/user_auth_info/%s/", dataID.String()), bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		mockService.On("UpdateUserAuthInfo", mock.Anything, mock.Anything, dataID, "testName", "testLogin", "testPassword", mock.Anything).Return(&models.UserAuthInfo{}, nil).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		invalidJSON := `{"name":"testName", "login": "testLogin", "password"`
		dataID := uuid.New()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/user_auth_info/%s/", dataID.String()), bytes.NewBufferString(invalidJSON))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "unexpected EOF")
	})

	t.Run("Missing Fields", func(t *testing.T) {
		dataID := uuid.New()
		requestBody, _ := json.Marshal(gin.H{
			"name":  "testName",
			"login": "testLogin",
			// "password" field is missing
			"metadata": gin.H{},
		})

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/user_auth_info/%s/", dataID.String()), bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
	t.Run("NotFound", func(t *testing.T) {
		dataID := uuid.New()
		requestBody, _ := json.Marshal(gin.H{
			"name":     "testName",
			"login":    "testLogin",
			"password": "testPassword",
			"metadata": gin.H{},
		})
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/user_auth_info/%s/", dataID.String()), bytes.NewBuffer(requestBody))

		rec := httptest.NewRecorder()

		mockService.On("UpdateUserAuthInfo", mock.Anything, mock.Anything, dataID, "testName", "testLogin", "testPassword", mock.Anything).Return(nil, httperror.New(nil, "not found", http.StatusNotFound)).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})
	t.Run("Invalid UUID", func(t *testing.T) {
		invalidID := "invalid-uuid"
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/user_auth_info/%s/", invalidID), nil)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid data id")
	})
}

func TestGetUserAuthInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockIUserDataService)
	handlers := NewUserDataHandlers(mockService)

	router := gin.Default()
	userID := uuid.New()
	authenticatedGroup := router.Group("/", outmiddlewares.NewAuthMiddlewareForTest(userID))
	authenticatedGroup.GET("/user_auth_info/:id/", handlers.GetUserAuthInfo)

	t.Run("Success", func(t *testing.T) {
		dataID := uuid.New()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/user_auth_info/%s/", dataID.String()), nil)

		rec := httptest.NewRecorder()

		mockService.On("GetUserAuthInfo", mock.Anything, dataID, userID).Return(&models.UserAuthInfo{}, nil).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		dataID := uuid.New()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/user_auth_info/%s/", dataID.String()), nil)

		rec := httptest.NewRecorder()

		mockService.On("GetUserAuthInfo", mock.Anything, mock.Anything, mock.Anything).Return(nil, httperror.New(nil, "not found", http.StatusNotFound)).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("BadRequestInvalidID", func(t *testing.T) {
		invalidID := "invalid-uuid"
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/user_auth_info/%s/", invalidID), nil)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid data id")
	})

	t.Run("InternalServerError", func(t *testing.T) {
		dataID := uuid.New()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/user_auth_info/%s/", dataID.String()), nil)

		rec := httptest.NewRecorder()

		mockService.On("GetUserAuthInfo", mock.Anything, dataID, userID).Return(nil, fmt.Errorf("internal error")).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "internal error")
		mockService.AssertExpectations(t)
	})
}

func TestInsertUserTextData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockIUserDataService)
	handlers := NewUserDataHandlers(mockService)

	router := gin.Default()
	userID := uuid.New()
	authenticatedGroup := router.Group("/", outmiddlewares.NewAuthMiddlewareForTest(userID))
	authenticatedGroup.POST("/user_text_data/", handlers.InsertUserTextData)

	t.Run("Success", func(t *testing.T) {
		requestBody, _ := json.Marshal(gin.H{
			"name":      "testName",
			"text_data": "testText",
			"metadata":  gin.H{},
		})

		req, _ := http.NewRequest(http.MethodPost, "/user_text_data/", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		mockService.On("InsertUserTextData", mock.Anything, mock.Anything, "testName", "testText", mock.Anything).Return(&models.UserTextData{}, nil).Once()

		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		invalidJSON := "{invalid-json}"
		req, _ := http.NewRequest(http.MethodPost, "/user_text_data/", bytes.NewBufferString(invalidJSON))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid character")
	})

	t.Run("Missing Fields", func(t *testing.T) {
		requestBody, _ := json.Marshal(gin.H{
			"name": "testName",
			// "text_data" field is missing
			"metadata": gin.H{},
		})

		req, _ := http.NewRequest(http.MethodPost, "/user_text_data/", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "name and text_data are required")
	})

	t.Run("Service Error", func(t *testing.T) {
		requestBody, _ := json.Marshal(gin.H{
			"name":      "testName",
			"text_data": "testText",
			"metadata":  gin.H{},
		})

		req, _ := http.NewRequest(http.MethodPost, "/user_text_data/", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		mockService.On("InsertUserTextData", mock.Anything, mock.Anything, "testName", "testText", mock.Anything).Return(nil, httperror.New(nil, "internal error", http.StatusInternalServerError)).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "internal error")
	})
}

func TestGetUserTextData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockIUserDataService)
	handlers := NewUserDataHandlers(mockService)

	router := gin.Default()
	userID := uuid.New()
	authenticatedGroup := router.Group("/", outmiddlewares.NewAuthMiddlewareForTest(userID))
	authenticatedGroup.GET("/user_text_data/:id/", handlers.GetUserTextData)

	t.Run("Success", func(t *testing.T) {
		dataID := uuid.New()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/user_text_data/%s/", dataID.String()), nil)

		rec := httptest.NewRecorder()

		mockService.On("GetUserTextData", mock.Anything, dataID, userID).Return(&models.UserTextData{}, nil).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("BadRequestInvalidID", func(t *testing.T) {
		// Тестируем случай с некорректным UUID
		req, _ := http.NewRequest(http.MethodGet, "/user_text_data/invalid-uuid/", nil)

		rec := httptest.NewRecorder()

		// Обработка запроса
		router.ServeHTTP(rec, req)

		// Проверка, что возвращается ошибка 400 (Invalid data id)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"detail":"Invalid data id"}`, rec.Body.String())
	})

	t.Run("NotFound", func(t *testing.T) {
		dataID := uuid.New()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/user_text_data/%s/", dataID.String()), nil)

		rec := httptest.NewRecorder()

		mockService.On("GetUserTextData", mock.Anything, mock.Anything, mock.Anything).Return(nil, httperror.New(nil, "not found", http.StatusNotFound)).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUpdateUserTextData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockIUserDataService)
	handlers := NewUserDataHandlers(mockService)

	router := gin.Default()
	userID := uuid.New()
	authenticatedGroup := router.Group("/", outmiddlewares.NewAuthMiddlewareForTest(userID))
	authenticatedGroup.PUT("/user_text_data/:id/", handlers.UpdateUserTextData)

	t.Run("Success", func(t *testing.T) {
		dataID := uuid.New()
		requestBody, _ := json.Marshal(gin.H{
			"name":      "testName",
			"text_data": "testText",
			"metadata":  gin.H{},
		})

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/user_text_data/%s/", dataID.String()), bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		mockService.On("UpdateUserTextData", mock.Anything, userID, dataID, "testName", "testText", mock.Anything).Return(&models.UserTextData{}, nil).Once()

		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, "/user_text_data/invalid-uuid/", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid data id")
	})

	t.Run("Service Error", func(t *testing.T) {
		dataID := uuid.New()
		requestBody, _ := json.Marshal(gin.H{
			"name":      "testName",
			"text_data": "testText",
			"metadata":  gin.H{},
		})

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/user_text_data/%s/", dataID.String()), bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		mockService.On("UpdateUserTextData", mock.Anything, userID, dataID, "testName", "testText", mock.Anything).Return(nil, httperror.New(nil, "internal error", http.StatusInternalServerError)).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "internal error")
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		dataID := uuid.New()
		invalidJSON := "{invalid-json}"
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/user_text_data/%s/", dataID.String()), bytes.NewBufferString(invalidJSON))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid character")
	})

	t.Run("Missing Fields", func(t *testing.T) {
		dataID := uuid.New()
		requestBody, _ := json.Marshal(gin.H{
			"name": "testName",
			// "text_data" field is missing
			"metadata": gin.H{},
		})

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/user_text_data/%s/", dataID.String()), bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "name and text_data are required")
	})
}

func TestDeleteUserTextData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockIUserDataService)
	handlers := NewUserDataHandlers(mockService)

	router := gin.Default()
	userID := uuid.New()
	authenticatedGroup := router.Group("/", outmiddlewares.NewAuthMiddlewareForTest(userID))
	authenticatedGroup.DELETE("/user_text_data/:id/", handlers.DeleteUserTextData)

	t.Run("Success", func(t *testing.T) {
		dataID := uuid.New()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/user_text_data/%s/", dataID.String()), nil)

		rec := httptest.NewRecorder()

		mockService.On("DeleteUserTextData", mock.Anything, dataID, userID).Return(nil).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		mockService.AssertExpectations(t)
	})
	t.Run("BadRequestInvalidID", func(t *testing.T) {
		// Тестируем случай с некорректным UUID
		req, _ := http.NewRequest(http.MethodDelete, "/user_text_data/invalid-uuid/", nil)

		rec := httptest.NewRecorder()

		// Обработка запроса
		router.ServeHTTP(rec, req)

		// Проверка, что возвращается ошибка 400 (Invalid data id)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"detail":"Invalid data id"}`, rec.Body.String())
	})

	t.Run("NotFound", func(t *testing.T) {
		dataID := uuid.New()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/user_text_data/%s/", dataID.String()), nil)

		rec := httptest.NewRecorder()

		mockService.On("DeleteUserTextData", mock.Anything, mock.Anything, mock.Anything).Return(httperror.New(nil, "not found", http.StatusNotFound)).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestInsertUserFileData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockIUserDataService)
	handlers := NewUserDataHandlers(mockService)

	router := gin.Default()
	userID := uuid.New()
	authenticatedGroup := router.Group("/", outmiddlewares.NewAuthMiddlewareForTest(userID))
	authenticatedGroup.POST("/user_file_data/", handlers.InsertUserFileData)

	t.Run("Success", func(t *testing.T) {
		// Создаем мультипарт-форму с файлом
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		fileWriter, err := writer.CreateFormFile("file", "testFile.txt")
		if err != nil {
			t.Fatalf("Failed to create form file: %v", err)
		}
		_, err = fileWriter.Write([]byte("This is a test file"))
		if err != nil {
			t.Fatalf("Failed to write to form file: %v", err)
		}
		writer.Close()

		req, _ := http.NewRequest(http.MethodPost, "/user_file_data/", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()
		mockService.On("InsertUserFileData", mock.Anything, userID, mock.Anything, mock.Anything).Return(&models.UserFileData{}, nil).Once()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("BadRequestMissingFile", func(t *testing.T) {
		// Тестируем случай без файла
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.Close() // Без файла

		req, _ := http.NewRequest(http.MethodPost, "/user_file_data/", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		rec := httptest.NewRecorder()

		// Проверка, что возвращается ошибка 400 (Bad Request) при отсутствии файла
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"detail":"http: no such file"}`, rec.Body.String())
	})

	t.Run("ErrorSavingFile", func(t *testing.T) {
		// Мокируем ошибку при сохранении файла
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		fileWriter, err := writer.CreateFormFile("file", "testFile.txt")
		if err != nil {
			t.Fatalf("Failed to create form file: %v", err)
		}
		_, err = fileWriter.Write([]byte("This is a test file"))
		if err != nil {
			t.Fatalf("Failed to write to form file: %v", err)
		}
		writer.Close()

		req, _ := http.NewRequest(http.MethodPost, "/user_file_data/", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		rec := httptest.NewRecorder()

		// Мокируем ошибку сохранения файла
		mockService.On("InsertUserFileData", mock.Anything, userID, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("failed to save file")).Once()

		// Проверка на ошибку
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t, `{"detail":"failed to save file"}`, rec.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("ErrorCreatingDirectory", func(t *testing.T) {
		// Мокируем ошибку при создании директории
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		fileWriter, err := writer.CreateFormFile("file", "testFile.txt")
		if err != nil {
			t.Fatalf("Failed to create form file: %v", err)
		}
		_, err = fileWriter.Write([]byte("This is a test file"))
		if err != nil {
			t.Fatalf("Failed to write to form file: %v", err)
		}
		writer.Close()

		req, _ := http.NewRequest(http.MethodPost, "/user_file_data/", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		rec := httptest.NewRecorder()

		// Мокируем ошибку при создании директории
		mockService.On("InsertUserFileData", mock.Anything, userID, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("failed to create directory")).Once()

		// Проверка на ошибку
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t, `{"detail":"failed to create directory"}`, rec.Body.String())
		mockService.AssertExpectations(t)
	})
}

func TestGetUserFileData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockIUserDataService)
	handlers := NewUserDataHandlers(mockService)

	router := gin.Default()
	userID := uuid.New()
	authenticatedGroup := router.Group("/", outmiddlewares.NewAuthMiddlewareForTest(userID))
	authenticatedGroup.GET("/user_file_data/:id/", handlers.GetUserFileData)

	t.Run("Success", func(t *testing.T) {
		dataID := uuid.New()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/user_file_data/%s/", dataID.String()), nil)

		rec := httptest.NewRecorder()

		mockService.On("GetUserFileData", mock.Anything, dataID, userID).Return(&models.UserFileData{}, nil).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		dataID := uuid.New()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/user_file_data/%s/", dataID.String()), nil)

		rec := httptest.NewRecorder()

		mockService.On("GetUserFileData", mock.Anything, dataID, userID).Return(&models.UserFileData{}, httperror.New(nil, "not found", http.StatusNotFound)).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})
	t.Run("BadRequestInvalidID", func(t *testing.T) {
		// Тестируем случай с некорректным UUID
		req, _ := http.NewRequest(http.MethodGet, "/user_file_data/invalid-uuid/", nil)

		rec := httptest.NewRecorder()

		// Обработка запроса
		router.ServeHTTP(rec, req)

		// Проверка, что возвращается ошибка 400 (Invalid data id)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"detail":"Invalid data id"}`, rec.Body.String())
	})
}

func TestDeleteUserFileData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockIUserDataService)
	handlers := NewUserDataHandlers(mockService)

	router := gin.Default()
	userID := uuid.New()
	authenticatedGroup := router.Group("/", outmiddlewares.NewAuthMiddlewareForTest(userID))
	authenticatedGroup.DELETE("/user_file_data/:id/", handlers.DeleteUserFileData)

	t.Run("Success", func(t *testing.T) {
		// Успешный случай удаления
		dataID := uuid.New()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/user_file_data/%s/", dataID.String()), nil)

		rec := httptest.NewRecorder()

		// Мокируем успешное удаление
		mockService.On("DeleteUserFileData", mock.Anything, dataID, userID).Return(nil).Once()

		// Обработка запроса
		router.ServeHTTP(rec, req)

		// Проверка, что ответ статус 204 (No Content)
		assert.Equal(t, http.StatusNoContent, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		// Удаление не найденного файла
		dataID := uuid.New()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/user_file_data/%s/", dataID.String()), nil)

		rec := httptest.NewRecorder()

		// Мокируем ошибку "not found"
		mockService.On("DeleteUserFileData", mock.Anything, dataID, userID).Return(httperror.New(nil, "not found", http.StatusNotFound)).Once()

		// Обработка запроса
		router.ServeHTTP(rec, req)

		// Проверка, что статус ошибки - 404
		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("BadRequestInvalidID", func(t *testing.T) {
		// Тестируем случай с некорректным UUID
		req, _ := http.NewRequest(http.MethodDelete, "/user_file_data/invalid-uuid/", nil)

		rec := httptest.NewRecorder()

		// Обработка запроса
		router.ServeHTTP(rec, req)

		// Проверка, что возвращается ошибка 400 (Invalid data id)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"detail":"Invalid data id"}`, rec.Body.String())
	})
}

func TestUpdateUserFileData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockIUserDataService)
	handlers := NewUserDataHandlers(mockService)

	router := gin.Default()
	userID := uuid.New()
	authenticatedGroup := router.Group("/", outmiddlewares.NewAuthMiddlewareForTest(userID))
	authenticatedGroup.PUT("/user_file_data/:id/", handlers.UpdateUserFileData)

	t.Run("Success", func(t *testing.T) {
		requestBody, _ := json.Marshal(gin.H{
			"name":     "testName",
			"metadata": gin.H{},
		})
		dataID := uuid.New()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/user_file_data/%s/", dataID.String()), bytes.NewBuffer(requestBody))

		rec := httptest.NewRecorder()

		// Настройка мока для успешного вызова
		mockService.On("UpdateUserFileData", mock.Anything, userID, dataID, mock.Anything, mock.Anything).Return(&models.UserFileData{}, nil).Once()

		router.ServeHTTP(rec, req)

		// Проверка успешного кода ответа
		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		dataID := uuid.New()
		requestBody, _ := json.Marshal(gin.H{
			"name":     "testName",
			"metadata": gin.H{},
		})
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/user_file_data/%s/", dataID.String()), bytes.NewBuffer(requestBody))

		rec := httptest.NewRecorder()

		// Настройка мока для ошибки "not found"
		mockService.On("UpdateUserFileData", mock.Anything, userID, dataID, mock.Anything, mock.Anything).Return(nil, httperror.New(nil, "not found", http.StatusNotFound)).Once()

		router.ServeHTTP(rec, req)

		// Проверка кода ошибки
		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		// Тестирование некорректного ID
		req, _ := http.NewRequest(http.MethodPut, "/user_file_data/invalid-uuid/", nil)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		// Проверка кода ошибки "BadRequest"
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"detail":"Invalid data id"}`, rec.Body.String())
	})

	t.Run("BadRequestMissingName", func(t *testing.T) {
		// Тестирование отсутствующего поля "name"
		requestBody, _ := json.Marshal(gin.H{
			"metadata": gin.H{},
		})
		dataID := uuid.New()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/user_file_data/%s/", dataID.String()), bytes.NewBuffer(requestBody))

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		// Проверка ошибки для отсутствующего поля "name"
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"detail":"name is required"}`, rec.Body.String())
	})

	t.Run("BadRequestInvalidJson", func(t *testing.T) {
		dataID := uuid.New()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/user_file_data/%s/", dataID.String()), bytes.NewBufferString("invalid json"))

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		// Проверка ошибки для некорректного JSON
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"detail":"invalid character 'i' looking for beginning of value"}`, rec.Body.String())
	})
}

func TestDownloadUserFile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockIUserDataService)
	handlers := NewUserDataHandlers(mockService)
	router := gin.Default()
	userID := uuid.New()
	tempFile, err := os.CreateTemp("", "testfile-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	testFileContent := "This is a test file"
	_, err = tempFile.WriteString(testFileContent)
	if err != nil {
		t.Fatal(err)
	}
	tempFile.Close()

	authenticatedGroup := router.Group("/", outmiddlewares.NewAuthMiddlewareForTest(userID))
	authenticatedGroup.GET("/user_file_data/:id/download/", handlers.DownloadUserFile)
	t.Run("Invalid DataID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/user_file_data/invalid-id/download/", nil)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"detail":"Invalid data id"}`, rec.Body.String())
	})
	t.Run("NotFound", func(t *testing.T) {
		dataID := uuid.New()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/user_file_data/%s/download/", dataID.String()), nil)

		rec := httptest.NewRecorder()

		mockService.On("GetUserFileData", mock.Anything, dataID, userID).Return(&models.UserFileData{}, httperror.New(nil, "not found", http.StatusNotFound)).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})
	t.Run("Success", func(t *testing.T) {
		dataID := uuid.New()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/user_file_data/%s/download/", dataID.String()), nil)

		rec := httptest.NewRecorder()

		mockService.On("GetUserFileData", mock.Anything, dataID, userID).Return(&models.UserFileData{PathToFile: tempFile.Name()}, nil).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, rec.Body.String(), testFileContent)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestInsertUserBankCard(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockIUserDataService)
	handlers := NewUserDataHandlers(mockService)

	router := gin.Default()
	userID := uuid.New()
	authenticatedGroup := router.Group("/", outmiddlewares.NewAuthMiddlewareForTest(userID))
	authenticatedGroup.POST("/user_bank_card/", handlers.InsertUserBankCard)

	t.Run("Success", func(t *testing.T) {
		requestBody, _ := json.Marshal(gin.H{
			"name":        "testName",
			"number":      "testNumber",
			"card_holder": "testCardHolder",
			"expire_date": "testExpireDate",
			"csc":         "testCSC",
			"metadata":    gin.H{},
		})

		req, _ := http.NewRequest(http.MethodPost, "/user_bank_card/", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		mockService.On("InsertUserBankCard", mock.Anything, userID, "testName", "testNumber", "testCardHolder", "testExpireDate", "testCSC", mock.Anything).Return(&models.UserBankCard{}, nil).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		invalidJSON := `{"name":"testName","number":12345` // Некорректный JSON

		req, _ := http.NewRequest(http.MethodPost, "/user_bank_card/", bytes.NewBufferString(invalidJSON))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "unexpected EOF")
	})

	t.Run("Missing Required Fields", func(t *testing.T) {
		requestBody, _ := json.Marshal(gin.H{
			"name": "testName", // Отсутствуют остальные обязательные поля
		})

		req, _ := http.NewRequest(http.MethodPost, "/user_bank_card/", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"detail":"name,number,card_holder,expire_date and csc are required"}`, rec.Body.String())
	})

	t.Run("Service Error", func(t *testing.T) {
		requestBody, _ := json.Marshal(gin.H{
			"name":        "testName",
			"number":      "testNumber",
			"card_holder": "testCardHolder",
			"expire_date": "testExpireDate",
			"csc":         "testCSC",
			"metadata":    gin.H{},
		})

		req, _ := http.NewRequest(http.MethodPost, "/user_bank_card/", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		mockService.On("InsertUserBankCard", mock.Anything, userID, "testName", "testNumber", "testCardHolder", "testExpireDate", "testCSC", mock.Anything).Return(nil, errors.New("database error")).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t, `{"detail":"database error"}`, rec.Body.String())
	})
}

func TestGetUserBankCard(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockIUserDataService)
	handlers := NewUserDataHandlers(mockService)

	router := gin.Default()
	userID := uuid.New()
	authenticatedGroup := router.Group("/", outmiddlewares.NewAuthMiddlewareForTest(userID))
	authenticatedGroup.GET("/user_bank_card/:id/", handlers.GetUserBankCard)

	t.Run("Invalid DataID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/user_bank_card/invalid-id/", nil)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"detail":"Invalid data id"}`, rec.Body.String())
	})

	t.Run("Success", func(t *testing.T) {
		dataID := uuid.New()
		expectedCard := &models.UserBankCard{
			BaseUserData: models.BaseUserData{
				ID:       dataID,
				UserID:   userID,
				Name:     "testName",
				Metadata: map[string]interface{}{},
			},

			Number:     "testNumber",
			CardHolder: "testCardHolder",
			ExpireDate: "testExpireDate",
			CSC:        "testCSC",
		}

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/user_bank_card/%s/", dataID.String()), nil)

		rec := httptest.NewRecorder()

		mockService.On("GetUserBankCard", mock.Anything, dataID, userID).Return(expectedCard, nil).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var actualCard models.UserBankCard
		err := json.Unmarshal(rec.Body.Bytes(), &actualCard)
		assert.NoError(t, err)
		assert.Equal(t, expectedCard, &actualCard)

		mockService.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		dataID := uuid.New()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/user_bank_card/%s/", dataID.String()), nil)

		rec := httptest.NewRecorder()

		mockService.On("GetUserBankCard", mock.Anything, dataID, userID).Return(nil, httperror.New(nil, "not found", http.StatusNotFound)).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.JSONEq(t, `{"detail":"not found"}`, rec.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		dataID := uuid.New()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/user_bank_card/%s/", dataID.String()), nil)

		rec := httptest.NewRecorder()

		mockService.On("GetUserBankCard", mock.Anything, dataID, userID).Return(nil, errors.New("database error")).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t, `{"detail":"database error"}`, rec.Body.String())

		mockService.AssertExpectations(t)
	})
}

func TestDeleteUserBankCard(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockIUserDataService)
	handlers := NewUserDataHandlers(mockService)

	router := gin.Default()
	userID := uuid.New()
	authenticatedGroup := router.Group("/", outmiddlewares.NewAuthMiddlewareForTest(userID))
	authenticatedGroup.DELETE("/user_bank_card/:id/", handlers.DeleteUserBankCard)

	t.Run("Invalid DataID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/user_bank_card/invalid-id/", nil)

		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"detail":"Invalid data id"}`, rec.Body.String())
	})

	t.Run("Success", func(t *testing.T) {
		dataID := uuid.New()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/user_bank_card/%s/", dataID.String()), nil)

		rec := httptest.NewRecorder()

		mockService.On("DeleteUserBankCard", mock.Anything, dataID, userID).Return(nil).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		dataID := uuid.New()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/user_bank_card/%s/", dataID.String()), nil)

		rec := httptest.NewRecorder()

		mockService.On("DeleteUserBankCard", mock.Anything, dataID, userID).Return(httperror.New(nil, "not found", http.StatusNotFound)).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.JSONEq(t, `{"detail":"not found"}`, rec.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		dataID := uuid.New()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/user_bank_card/%s/", dataID.String()), nil)

		rec := httptest.NewRecorder()

		mockService.On("DeleteUserBankCard", mock.Anything, dataID, userID).Return(errors.New("database error")).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t, `{"detail":"database error"}`, rec.Body.String())
		mockService.AssertExpectations(t)
	})
}

func TestUpdateUserBankCard(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockIUserDataService)
	handlers := NewUserDataHandlers(mockService)

	router := gin.Default()
	userID := uuid.New()
	authenticatedGroup := router.Group("/", outmiddlewares.NewAuthMiddlewareForTest(userID))
	authenticatedGroup.PUT("/user_bank_card/:id/", handlers.UpdateUserBankCard)

	t.Run("Invalid DataID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, "/user_bank_card/invalid-id/", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"detail":"Invalid data id"}`, rec.Body.String())
	})

	t.Run("Invalid JSON Body", func(t *testing.T) {
		dataID := uuid.New()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/user_bank_card/%s/", dataID.String()), bytes.NewBufferString("{invalid json"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "detail")
	})

	t.Run("Missing Required Fields", func(t *testing.T) {
		dataID := uuid.New()
		requestBody, _ := json.Marshal(gin.H{"name": "testName"}) // Остальные поля отсутствуют
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/user_bank_card/%s/", dataID.String()), bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"detail":"name,number,card_holder,expire_date and csc are required"}`, rec.Body.String())
	})

	t.Run("Service Error", func(t *testing.T) {
		dataID := uuid.New()
		requestBody, _ := json.Marshal(gin.H{
			"name":        "testName",
			"number":      "testNumber",
			"card_holder": "testCardHolder",
			"expire_date": "testExpireDate",
			"csc":         "testCSC",
			"metadata":    gin.H{},
		})
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/user_bank_card/%s/", dataID.String()), bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mockService.On(
			"UpdateUserBankCard",
			mock.Anything,
			userID,
			dataID,
			"testName",
			"testNumber",
			"testCardHolder",
			"testExpireDate",
			"testCSC",
			mock.Anything,
		).Return(nil, errors.New("database error")).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t, `{"detail":"database error"}`, rec.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		dataID := uuid.New()
		requestBody, _ := json.Marshal(gin.H{
			"name":        "testName",
			"number":      "testNumber",
			"card_holder": "testCardHolder",
			"expire_date": "testExpireDate",
			"csc":         "testCSC",
			"metadata":    gin.H{},
		})
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/user_bank_card/%s/", dataID.String()), bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mockService.On(
			"UpdateUserBankCard",
			mock.Anything,
			userID,
			dataID,
			"testName",
			"testNumber",
			"testCardHolder",
			"testExpireDate",
			"testCSC",
			mock.Anything,
		).Return(&models.UserBankCard{}, nil).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})
}
