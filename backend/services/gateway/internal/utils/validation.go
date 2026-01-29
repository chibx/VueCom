package utils

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

var _validator = validator.New(validator.WithRequiredStructEnabled())

func Validator() *validator.Validate {
	return _validator
}

func TagNameFunc(fld reflect.StructField) string {
	name := fld.Tag.Get("name")
	if name == "" {
		return fld.Name
	}
	return name
}
