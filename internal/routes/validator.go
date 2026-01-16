package routes

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewValidator() *CustomValidator {
	v := validator.New()

	// register custom validations ONCE
	v.RegisterValidation("strpwd", strongPassword)
	v.RegisterValidation("phone", phone)
	v.RegisterValidation("aadhar", AadharNumber)
	v.RegisterValidation("pan", PanNumber)

	return &CustomValidator{
		validator: v,
	}
}

func (cv *CustomValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}

// Strong password validator function
func strongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool

	for _, ch := range password {
		switch {
		case ch >= 'A' && ch <= 'Z':
			hasUpper = true
		case ch >= 'a' && ch <= 'z':
			hasLower = true
		case ch >= '0' && ch <= '9':
			hasNumber = true
		default:
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// Phone number validator function
func phone(fl validator.FieldLevel) bool {
	var phoneRegex = regexp.MustCompile(`^[0-9]{10}$`)
	phone := fl.Field().String()
	return phoneRegex.MatchString(phone)
}

func AadharNumber(fl validator.FieldLevel) bool {
	aadhar := fl.Field().String()
	regex := regexp.MustCompile(`^[0-9]{12}$`)
	return regex.MatchString(aadhar)
}

func PanNumber(fl validator.FieldLevel) bool {
	pan := fl.Field().String()
	regex := regexp.MustCompile(`^[A-Z]{5}[0-9]{4}[A-Z]{1}$`)
	return regex.MatchString(pan)
}
