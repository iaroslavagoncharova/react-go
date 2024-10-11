package validation

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

// Initialize the validator
func init() {
    validate = validator.New()
}

// ValidateStruct validates the given struct
func ValidateStruct(s interface{}) error {
    return validate.Struct(s)
}