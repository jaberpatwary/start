package handler

import (
	"net/http"

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

	user, token, err := h.authUsecase.Register(c.Request().Context(), req.Name, req.Email, req.Password, req.Phone)
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
		HttpOnly: false,
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

	user, token, err := h.authUsecase.Login(c.Request().Context(), req.Email, req.Password)
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
		HttpOnly: false,
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

	user, err := h.authUsecase.UpdateProfile(c.Request().Context(), userID, req.Name, req.Phone, req.Avatar)
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
