package server

import (
	"regexp"

	valid "github.com/go-playground/validator/v10"
)

var (
	retailRegex      = regexp.MustCompile(`^[\w\s\-&]+$`)
	totalRegex       = regexp.MustCompile(`^\d+\.\d{2}$`)
	descriptionRegex = regexp.MustCompile(`^[\w\s\-]+$`)
	priceRegex       = regexp.MustCompile(`^\d+\.\d{2}$"`)
)

func newReciptValidator() *valid.Validate {

	validator := valid.New()

	validator.RegisterValidation("retailer", validateRetail, true)
	validator.RegisterValidation("total", validateTotal, true)
	validator.RegisterValidation("description", validateDescription, true)
	validator.RegisterValidation("price", validatePrice, true)

	return validator
}

func validateRetail(fl valid.FieldLevel) bool {
	return retailRegex.MatchString(fl.Field().String())
}

func validateTotal(fl valid.FieldLevel) bool {
	return totalRegex.MatchString(fl.Field().String())
}

func validateDescription(fl valid.FieldLevel) bool {
	return descriptionRegex.MatchString(fl.Field().String())
}

func validatePrice(fl valid.FieldLevel) bool {
	return priceRegex.MatchString(fl.Field().String())
}
