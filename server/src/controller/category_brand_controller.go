package controller

import (
	"app/src/response"
	"app/src/service"
	"app/src/validation"

	"github.com/gofiber/fiber/v2"
)

type CategoryBrandController struct {
	service service.CategoryBrandService
}

func NewCategoryBrandController(service service.CategoryBrandService) *CategoryBrandController {
	return &CategoryBrandController{service: service}
}

func (cbc *CategoryBrandController) GetCategories(c *fiber.Ctx) error {
	categories, err := cbc.service.GetAllCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorDetails{
			Code:    fiber.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":       fiber.StatusOK,
		"status":     "success",
		"categories": categories,
	})
}

func (cbc *CategoryBrandController) CreateCategory(c *fiber.Ctx) error {
	var req validation.CategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	category, err := cbc.service.CreateCategory(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code":     fiber.StatusCreated,
		"status":   "success",
		"category": category,
	})
}

func (cbc *CategoryBrandController) GetBrands(c *fiber.Ctx) error {
	brands, err := cbc.service.GetAllBrands()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorDetails{
			Code:    fiber.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":   fiber.StatusOK,
		"status": "success",
		"brands": brands,
	})
}

func (cbc *CategoryBrandController) CreateBrand(c *fiber.Ctx) error {
	var req validation.BrandRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	brand, err := cbc.service.CreateBrand(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code":   fiber.StatusCreated,
		"status": "success",
		"brand":  brand,
	})
}
