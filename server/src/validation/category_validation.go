package validation

type CategoryRequest struct {
	Name string `json:"name" validate:"required"`
	Icon string `json:"icon"`
}

type BrandRequest struct {
	Name string `json:"name" validate:"required"`
	Logo string `json:"logo"`
}
