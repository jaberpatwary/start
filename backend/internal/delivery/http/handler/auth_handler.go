package handler

import (
	"net/http"
	"net/mail"
	"strings"

	"github.com/jaberpatwary/startech/internal/domain"
	"github.com/jaberpatwary/startech/internal/usecase"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateProfileRequest struct {
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Avatar string `json:"avatar"`
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	cleanName := strings.TrimSpace(req.Name)
	cleanEmail := strings.ToLower(strings.TrimSpace(req.Email))
	cleanPhone := strings.TrimSpace(req.Phone)

	if cleanName == "" || len(cleanName) < 2 || len(cleanName) > 100 {
		return echo.NewHTTPError(http.StatusBadRequest, "Name must be between 2 and 100 characters")
	}

	if cleanEmail == "" || len(cleanEmail) > 100 {
		return echo.NewHTTPError(http.StatusBadRequest, "Valid email address is required")
	}
	if _, err := mail.ParseAddress(cleanEmail); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid email address format")
	}

	if len(req.Password) < 6 {
		return echo.NewHTTPError(http.StatusBadRequest, "Password must be at least 6 characters long")
	}
	if len(req.Password) > 72 {
		return echo.NewHTTPError(http.StatusBadRequest, "Password cannot exceed 72 characters")
	}

	user, token, err := h.authUsecase.Register(c.Request().Context(), cleanName, cleanEmail, req.Password, cleanPhone)
	if err != nil {
		if err == domain.ErrAlreadyExists {
			return echo.NewHTTPError(http.StatusConflict, "User with this email already exists")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.SetCookie(&http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	return c.JSON(http.StatusCreated, echo.Map{
		"code":    201,
		"status":  "success",
		"message": "Registration successful",
		"token":   token,
		"user":    sanitizeUser(user),
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	cleanEmail := strings.ToLower(strings.TrimSpace(req.Email))
	if cleanEmail == "" || req.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Email and password are required")
	}
	if len(req.Password) > 72 {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid email or password")
	}

	user, token, err := h.authUsecase.Login(c.Request().Context(), cleanEmail, req.Password)
	if err != nil {
		if err == domain.ErrAccountBlocked {
			return echo.NewHTTPError(http.StatusForbidden, "Account is blocked")
		}
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid email or password")
	}

	c.SetCookie(&http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	return c.JSON(http.StatusOK, echo.Map{
		"code":    200,
		"status":  "success",
		"message": "Login successful",
		"token":   token,
		"user":    sanitizeUser(user),
	})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:   "auth_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	return c.JSON(http.StatusOK, echo.Map{
		"code":    200,
		"status":  "success",
		"message": "Logged out successfully",
	})
}

func (h *AuthHandler) GetMe(c echo.Context) error {
	userID := c.Get("user_id").(string)
	user, err := h.authUsecase.GetProfile(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}
	return c.JSON(http.StatusOK, echo.Map{
		"code":   200,
		"status": "success",
		"user":   sanitizeUser(user),
	})
}

func (h *AuthHandler) UpdateProfile(c echo.Context) error {
	userID := c.Get("user_id").(string)
	var req UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	cleanName := strings.TrimSpace(req.Name)
	if cleanName == "" || len(cleanName) < 2 || len(cleanName) > 100 {
		return echo.NewHTTPError(http.StatusBadRequest, "Name must be between 2 and 100 characters")
	}

	user, err := h.authUsecase.UpdateProfile(c.Request().Context(), userID, cleanName, strings.TrimSpace(req.Phone), strings.TrimSpace(req.Avatar))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update profile")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":    200,
		"status":  "success",
		"message": "Profile updated successfully",
		"user":    sanitizeUser(user),
	})
}

func (h *AuthHandler) AddAddress(c echo.Context) error {
	userID := c.Get("user_id").(string)
	var addr domain.Address
	if err := c.Bind(&addr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid address data")
	}

	if err := h.authUsecase.AddAddress(c.Request().Context(), userID, &addr); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to add address")
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"code":    201,
		"status":  "success",
		"message": "Address added",
		"address": addr,
	})
}

func (h *AuthHandler) DeleteAddress(c echo.Context) error {
	userID := c.Get("user_id").(string)
	addressID := c.Param("id")
	if err := h.authUsecase.DeleteAddress(c.Request().Context(), addressID, userID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete address")
	}
	return c.JSON(http.StatusOK, echo.Map{
		"code":    200,
		"status":  "success",
		"message": "Address deleted",
	})
}

func sanitizeUser(u *domain.User) echo.Map {
	return echo.Map{
		"id":         u.ID,
		"name":       u.Name,
		"email":      u.Email,
		"phone":      u.Phone,
		"role":       u.Role,
		"avatar":     u.Avatar,
		"is_blocked": u.IsBlocked,
		"addresses":  u.Addresses,
		"created_at": u.CreatedAt,
	}
}
