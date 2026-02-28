package server

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var ErrDBRecordNotFound = errors.New("Record Not Found!!!")

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ServerErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *ServerErr) Error() string {
	return e.Message
}

func NewServerErr(code int, message string) *ServerErr {
	return &ServerErr{Code: code, Message: message}
}

func ValErrToBag(err error) []ErrorDetail {
	var validationErr, ok = err.(validator.ValidationErrors)
	var errorBag = make([]ErrorDetail, 0)
	if !ok {
		return errorBag
	}
	for _, v := range validationErr {
		field := v.Field()
		message := v.Error()

		error := ErrorDetail{
			Field:   field,
			Message: message,
		}

		errorBag = append(errorBag, error)
	}

	return errorBag
}
