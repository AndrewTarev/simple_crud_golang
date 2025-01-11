package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"

	"simple_crud_go/internal/db/models"
	"simple_crud_go/internal/service/mocks"
)

func TestCreateUser_Success_ConvertData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)

	mockService.EXPECT().
		CreateUser(gomock.Any(), gomock.Any()).
		Return(1, nil)

	handler := Handler{services: mockService}

	r := gin.Default()
	r.POST("/user", handler.CreateUser)

	reqBody := `{"username":"testuser","email":"test@example.com","password":"testpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Декодируем JSON в map
	var actualResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)

	// Преобразуем данные в структуру
	var userIdResponse models.UserIdResponse
	data, _ := actualResponse["data"].(map[string]interface{})
	userIdResponse.ID = int(data["id"].(float64))

	// Сравниваем значения
	assert.Equal(t, models.UserIdResponse{ID: 1}, userIdResponse)
}

func TestCreateUser_UsernameAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	handler := NewHandler(mockService)

	// Подготавливаем моковый ответ
	mockPgError := &pgconn.PgError{
		Code:           "23505", // Код ошибки уникального ограничения
		ConstraintName: "users_username_key",
	}

	mockService.EXPECT().
		CreateUser(gomock.Any(), gomock.Any()).
		Return(0, mockPgError)

	router := gin.Default()
	router.POST("/user", handler.CreateUser)

	// Отправляем тестовый запрос
	w := httptest.NewRecorder()
	reqBody := `{"username":"testuser","email":"test@example.com","password":"password123"}`
	req, _ := http.NewRequest(http.MethodPost, "/user", strings.NewReader(reqBody))
	router.ServeHTTP(w, req)

	// Проверяем ответ
	assert.Equal(t, http.StatusForbidden, w.Code)

	expectedResponse := map[string]interface{}{
		"status": "failed",
		"error": map[string]interface{}{
			"message": "Username is already exist",
		},
	}

	var actualResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &actualResponse)

	assert.Equal(t, expectedResponse, actualResponse)
}

func TestCreateUser_EmailAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)

	// Подготавливаем моковый ответ
	mockPgError := &pgconn.PgError{
		Code:           "23505", // Код ошибки уникального ограничения
		ConstraintName: "users_email_key",
	}

	mockService.EXPECT().
		CreateUser(gomock.Any(), gomock.Any()).
		Return(0, mockPgError)

	handler := Handler{services: mockService}

	r := gin.Default()
	r.POST("/user", handler.CreateUser)

	reqBody := `{"username":"testuser","email":"test@example.com","password":"testpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)

	var actualResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)

	expectedResponse := map[string]interface{}{
		"status": "failed",
		"error": map[string]interface{}{
			"message": "Email is already exist",
		},
	}

	assert.Equal(t, expectedResponse, actualResponse)
}

func TestCreateUser_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)

	mockService.EXPECT().
		CreateUser(gomock.Any(), gomock.Any()).
		Return(0, fmt.Errorf("internal server error"))

	handler := Handler{services: mockService}

	r := gin.Default()
	r.POST("/user", handler.CreateUser)

	reqBody := `{"username":"testuser","email":"test@example.com","password":"testpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var actualResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)

	expectedResponse := map[string]interface{}{
		"status": "failed",
		"error": map[string]interface{}{
			"message": "Something went wrong",
		},
	}

	assert.Equal(t, expectedResponse, actualResponse)
}

func TestCreateUser_ValidationError_EmptyEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	handler := Handler{services: mockService}

	r := gin.Default()
	r.POST("/user", handler.CreateUser)

	// Некорректные данные (пустой email)
	reqBody := `{"username":"testuser","email":"","password":"testpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем, что возвращается ошибка с кодом 400 (Bad Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var actualResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)

	expectedResponse := map[string]interface{}{
		"status": "failed",
		"error": map[string]interface{}{
			"message": "Email is required",
		},
	}

	assert.Equal(t, expectedResponse, actualResponse)
}

func TestCreateUser_ValidationError_InvalidEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	handler := Handler{services: mockService}

	r := gin.Default()
	r.POST("/user", handler.CreateUser)

	// Некорректные данные (неправильный формат email)
	reqBody := `{"username":"testuser","email":"invalidemail","password":"testpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем, что возвращается ошибка с кодом 400 (Bad Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var actualResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)

	expectedResponse := map[string]interface{}{
		"status": "failed",
		"error": map[string]interface{}{
			"message": "Email must be a valid email",
		},
	}

	assert.Equal(t, expectedResponse, actualResponse)
}

func TestCreateUser_ValidationError_ShortPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	handler := Handler{services: mockService}

	r := gin.Default()
	r.POST("/user", handler.CreateUser)

	// Некорректные данные (слишком короткий пароль)
	reqBody := `{"username":"testuser","email":"test@example.com","password":"short"}`
	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем, что возвращается ошибка с кодом 400 (Bad Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var actualResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)

	expectedResponse := map[string]interface{}{
		"status": "failed",
		"error": map[string]interface{}{
			"message": "Password must be at least 8 characters",
		},
	}

	assert.Equal(t, expectedResponse, actualResponse)
}
