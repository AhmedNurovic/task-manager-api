package handler

import (
	"net/http"
	"strings"

	"github.com/ahmednurovic/task-manager-api/internal/model"
	"github.com/ahmednurovic/task-manager-api/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthService interface {
	Register(ctx *gin.Context, email, password string) (*model.User, error)
	Login(ctx *gin.Context, email, password string) (string, error)
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param input body RegisterRequest true "Register input"
// @Success 201 {object} model.User
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
func Register(authService service.AuthServicer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required,min=8"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			// Custom error for invalid email format
			if strings.Contains(err.Error(), "Email") && strings.Contains(err.Error(), "email") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		user, err := authService.Register(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, user)
	}
}

// Login godoc
// @Summary Login a user
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param input body LoginRequest true "Login input"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/login [post]
func Login(authService service.AuthServicer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token, err := authService.Login(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"password123"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

type TokenResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}
