package controller

import (
	"app/src/model"
	"app/src/response"
	"app/src/service"
	"app/src/validation"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	userService service.UserService
}

func NewAuthController(userService service.UserService) *AuthController {
	return &AuthController{userService: userService}
}

func (ac *AuthController) Register(c *fiber.Ctx) error {
	var req validation.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request payload",
			Errors:  err.Error(),
		})
	}

	user, token, err := ac.userService.Register(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "error",
			Message: err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HTTPOnly: false,
		SameSite: "Lax",
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code":    fiber.StatusCreated,
		"status":  "success",
		"message": "User registered successfully",
		"user":    user,
		"token":   token,
	})
}

func (ac *AuthController) Login(c *fiber.Ctx) error {
	var req validation.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request payload",
			Errors:  err.Error(),
		})
	}

	user, token, err := ac.userService.Login(&req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorDetails{
			Code:    fiber.StatusUnauthorized,
			Status:  "error",
			Message: err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HTTPOnly: false,
		SameSite: "Lax",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"status":  "success",
		"message": "Login successful",
		"user":    user,
		"token":   token,
	})
}

func (ac *AuthController) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HTTPOnly: false,
	})
	return c.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Logout successful",
	})
}

func (ac *AuthController) GetMe(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorDetails{
			Code:    fiber.StatusUnauthorized,
			Status:  "error",
			Message: "Unauthorized access",
		})
	}

	user, err := ac.userService.GetProfile(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorDetails{
			Code:    fiber.StatusNotFound,
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithData[*model.User]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "User profile fetched",
		Data:    user,
	})
}

func (ac *AuthController) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	var req validation.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request body",
		})
	}

	user, err := ac.userService.UpdateProfile(userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"status":  "success",
		"message": "Profile updated successfully",
		"user":    user,
	})
}

func (ac *AuthController) GetAllUsers(c *fiber.Ctx) error {
	query := &validation.QueryUser{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
		Search: c.Query("search", ""),
	}

	res, err := ac.userService.GetAll(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorDetails{
			Code:    fiber.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
