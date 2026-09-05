package domain

import "errors"

var (
	ErrNotFound           = errors.New("record not found")
	ErrAlreadyExists      = errors.New("record already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrAccountBlocked     = errors.New("account is blocked")
	ErrBadRequest         = errors.New("invalid request parameters")
	ErrOutOfStock         = errors.New("product is out of stock")
	ErrCouponInvalid      = errors.New("invalid or expired coupon")
	ErrInternalServer     = errors.New("internal server error")
)
