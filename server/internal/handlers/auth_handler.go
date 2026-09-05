package handlers

import (
	"net/http"

	"github.com/jaberpatwary/startech/internal/middleware"
	"github.com/jaberpatwary/startech/internal/models"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/google/uuid"
	"strings"
)

type AuthHandler struct {
	DB *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

type RegisterInput struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Phone    string `json:"phone"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (h *AuthHandler) Register(c echo.Context) error {
	var input RegisterInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	var existing models.User
	if err := h.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		return echo.NewHTTPError(http.StatusConflict, "User with this email already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash password")
	}

	user := models.User{
		ID:           uuid.NewString(),
		Name:         input.Name,
		Email:        input.Email,
		Phone:        input.Phone,
		PasswordHash: string(hash),
		Role:         models.RoleUser,
	}
	if err := h.DB.Create(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}

	token, err := middleware.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate token")
	}

	c.SetCookie(&http.Cookie{Name: "auth_token", Value: token, Path: "/", HttpOnly: false, SameSite: http.SameSiteLaxMode})

	return c.JSON(http.StatusCreated, echo.Map{
		"code": 201, "status": "success", "message": "Registration successful",
		"token": token, "user": sanitizeUser(user),
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var input LoginInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	var user models.User
	emailInput := strings.TrimSpace(input.Email)
	if err := h.DB.Where("LOWER(email) = LOWER(?) OR (role = 'ADMIN' AND LOWER(?) IN ('admin', 'admin@startech.local', 'admin@startech.com'))", emailInput, emailInput).First(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid email or password")
	}
	if user.IsBlocked {
		return echo.NewHTTPError(http.StatusForbidden, "Account is blocked")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid email or password")
	}

	token, err := middleware.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate token")
	}

	c.SetCookie(&http.Cookie{Name: "auth_token", Value: token, Path: "/", HttpOnly: false, SameSite: http.SameSiteLaxMode})

	return c.JSON(http.StatusOK, echo.Map{
		"code": 200, "status": "success", "message": "Login successful",
		"token": token, "user": sanitizeUser(user),
	})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{Name: "auth_token", Value: "", Path: "/", MaxAge: -1})
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "message": "Logged out successfully"})
}

func (h *AuthHandler) GetMe(c echo.Context) error {
	userID := c.Get("user_id").(string)
	var user models.User
	if err := h.DB.Preload("Addresses").Where("id = ?", userID).First(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "user": sanitizeUser(user)})
}

func (h *AuthHandler) UpdateProfile(c echo.Context) error {
	userID := c.Get("user_id").(string)
	var user models.User
	if err := h.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}
	var input struct {
		Name   string `json:"name"`
		Phone  string `json:"phone"`
		Avatar string `json:"avatar"`
	}
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Phone != "" {
		user.Phone = input.Phone
	}
	if input.Avatar != "" {
		user.Avatar = input.Avatar
	}
	h.DB.Save(&user)
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "user": sanitizeUser(user)})
}

func (h *AuthHandler) AddAddress(c echo.Context) error {
	userID := c.Get("user_id").(string)
	var input models.Address
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	input.ID = uuid.NewString()
	input.UserID = userID
	if input.IsDefault {
		h.DB.Model(&models.Address{}).Where("user_id = ?", userID).Update("is_default", false)
	}
	h.DB.Create(&input)
	return c.JSON(http.StatusCreated, echo.Map{"code": 201, "status": "success", "address": input})
}

func (h *AuthHandler) DeleteAddress(c echo.Context) error {
	userID := c.Get("user_id").(string)
	addressID := c.Param("id")
	h.DB.Where("id = ? AND user_id = ?", addressID, userID).Delete(&models.Address{})
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "message": "Address deleted"})
}

func sanitizeUser(u models.User) echo.Map {
	return echo.Map{
		"id": u.ID, "name": u.Name, "email": u.Email,
		"phone": u.Phone, "role": u.Role, "avatar": u.Avatar,
		"is_blocked": u.IsBlocked, "addresses": u.Addresses,
		"created_at": u.CreatedAt,
	}
}
