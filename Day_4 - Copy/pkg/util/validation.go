package util

import (
	"fmt"

	"github.com/go-playground/validator"
)

func Struct(s interface{}) error {
	err := validator.New().Struct(s)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			return NewErrFieldValidation(e)
		}
	}

	return nil
}

func NewErrFieldValidation(err validator.FieldError) error {
	if err.Param() == "" {
		return fmt.Errorf("%s: %v; format must be (%s)", err.Field(), err.Value(), err.ActualTag())
	}
	return fmt.Errorf("%s: %v; format must be (%s=%s)", err.Field(), err.Value(), err.ActualTag(), err.Param())
}
