package validator

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=64"`
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=200"`
	Description string  `json:"description" validate:"max=2000"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock" validate:"required,gte=0"`
	CategoryID  int     `json:"category_id" validate:"required,gt=0"`
}

type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

type OrderItemRequest struct {
	ProductID int `json:"product_id" validate:"required,gt=0"`
	Quantity  int `json:"quantity" validate:"required,gt=0"`
}

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}
