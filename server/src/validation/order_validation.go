package validation

type OrderItemRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,gt=0"`
}

type CreateOrderRequest struct {
	Items            []OrderItemRequest `json:"items" validate:"required,min=1"`
	ShippingName     string             `json:"shipping_name" validate:"required"`
	ShippingPhone    string             `json:"shipping_phone" validate:"required"`
	ShippingDivision string             `json:"shipping_division" validate:"required"`
	ShippingDistrict string             `json:"shipping_district" validate:"required"`
	ShippingThana    string             `json:"shipping_thana" validate:"required"`
	ShippingAddress  string             `json:"shipping_address" validate:"required"`
	ShippingPostal   string             `json:"shipping_postal"`
	PaymentMethod    string             `json:"payment_method" validate:"required"`
	CouponCode       string             `json:"coupon_code"`
	Note             string             `json:"note"`
}

type UpdateOrderStatusRequest struct {
	Status        string `json:"status"`
	PaymentStatus string `json:"payment_status"`
}
