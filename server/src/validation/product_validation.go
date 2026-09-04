package validation

type CreateProductRequest struct {
	Name             string   `json:"name" validate:"required"`
	SKU              string   `json:"sku" validate:"required"`
	CategoryID       string   `json:"category_id" validate:"required"`
	BrandID          string   `json:"brand_id" validate:"required"`
	Price            int      `json:"price" validate:"required,gt=0"`
	DiscountPrice    *int     `json:"discount_price"`
	Stock            int      `json:"stock"`
	ShortDescription string   `json:"short_description"`
	Description      string   `json:"description"`
	Images           []string `json:"images"`
	Featured         bool     `json:"featured"`
	Status           string   `json:"status"`
}

type UpdateProductRequest struct {
	Name             string   `json:"name"`
	Price            int      `json:"price"`
	DiscountPrice    *int     `json:"discount_price"`
	Stock            int      `json:"stock"`
	ShortDescription string   `json:"short_description"`
	Description      string   `json:"description"`
	Images           []string `json:"images"`
	Featured         bool     `json:"featured"`
	Status           string   `json:"status"`
}

type QueryProduct struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Search     string `json:"search"`
	CategoryID string `json:"category_id"`
	BrandID    string `json:"brand_id"`
	MinPrice   int    `json:"min_price"`
	MaxPrice   int    `json:"max_price"`
	Featured   bool   `json:"featured"`
	SortBy     string `json:"sort_by"`
}
