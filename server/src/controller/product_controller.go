package controller

import (
	"app/src/response"
	"app/src/service"
	"app/src/validation"

	"github.com/gofiber/fiber/v2"
)

type ProductController struct {
	productService service.ProductService
}

func NewProductController(productService service.ProductService) *ProductController {
	return &ProductController{productService: productService}
}

func (pc *ProductController) CreateProduct(c *fiber.Ctx) error {
	var req validation.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	product, err := pc.productService.CreateProduct(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code":    fiber.StatusCreated,
		"status":  "success",
		"message": "Product created successfully",
		"product": product,
	})
}

func (pc *ProductController) UpdateProduct(c *fiber.Ctx) error {
	productID := c.Params("id")
	var req validation.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request body",
		})
	}

	product, err := pc.productService.UpdateProduct(productID, &req)
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
		"message": "Product updated successfully",
		"product": product,
	})
}

func (pc *ProductController) DeleteProduct(c *fiber.Ctx) error {
	productID := c.Params("id")
	if err := pc.productService.DeleteProduct(productID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Product deleted successfully",
	})
}

func (pc *ProductController) GetBySlug(c *fiber.Ctx) error {
	productSlug := c.Params("slug")
	product, err := pc.productService.GetBySlug(productSlug)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorDetails{
			Code:    fiber.StatusNotFound,
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"status":  "success",
		"product": product,
	})
}

func (pc *ProductController) GetByID(c *fiber.Ctx) error {
	productID := c.Params("id")
	product, err := pc.productService.GetByID(productID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorDetails{
			Code:    fiber.StatusNotFound,
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"status":  "success",
		"product": product,
	})
}

func (pc *ProductController) GetAll(c *fiber.Ctx) error {
	query := &validation.QueryProduct{
		Page:       c.QueryInt("page", 1),
		Limit:      c.QueryInt("limit", 12),
		Search:     c.Query("search", ""),
		CategoryID: c.Query("category_id", ""),
		BrandID:    c.Query("brand_id", ""),
		MinPrice:   c.QueryInt("min_price", 0),
		MaxPrice:   c.QueryInt("max_price", 0),
		Featured:   c.QueryBool("featured", false),
		SortBy:     c.Query("sort_by", ""),
	}

	result, err := pc.productService.GetAll(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorDetails{
			Code:    fiber.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
