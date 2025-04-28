package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ahmednurovic/task-manager-api/internal/handler"
	"github.com/ahmednurovic/task-manager-api/internal/model"
	"github.com/ahmednurovic/task-manager-api/internal/service"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, email string, password string) (*model.User, error) {
	args := m.Called(ctx, email, password)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, email string, password string) (string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Error(1)
}

func TestRegisterHandler(t *testing.T) {
	t.Run("Successful Registration", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := `{"email":"test@example.com","password":"password123"}`
		req := httptest.NewRequest("POST", "/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		mockAuthService := new(MockAuthService)
		mockUser := &model.User{
			ID:    1,
			Email: "test@example.com",
		}
		mockAuthService.On("Register", mock.Anything, "test@example.com", "password123").
			Return(mockUser, nil)

		var authService service.AuthServicer = mockAuthService
		handler.Register(authService)(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, `{"id":1,"email":"test@example.com"}`, w.Body.String())
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Invalid Email Format", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := `{"email":"invalid-email","password":"password123"}`
		req := httptest.NewRequest("POST", "/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		var authService service.AuthServicer = new(MockAuthService)
		handler.Register(authService)(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid email")
	})
}

func TestLoginHandler(t *testing.T) {
	t.Run("Successful Login", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := `{"email":"test@example.com","password":"password123"}`
		req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		mockAuthService := new(MockAuthService)
		mockAuthService.On("Login", mock.Anything, "test@example.com", "password123").
			Return("valid.token", nil)

		// Convert mock to service.AuthServicer interface type
		var authService service.AuthServicer = mockAuthService
		handler.Login(authService)(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"token":"valid.token"}`, w.Body.String())
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := `{"email":"test@example.com","password":"wrong"}`
		req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		mockAuthService := new(MockAuthService)
		mockAuthService.On("Login", mock.Anything, "test@example.com", "wrong").
			Return("", service.ErrInvalidCredentials)

		var authService service.AuthServicer = mockAuthService
		handler.Login(authService)(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "invalid credentials")
		mockAuthService.AssertExpectations(t)
	})
}
