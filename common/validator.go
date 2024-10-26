package common

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validateSingleton = validator.New()

type errorResponse struct {
	Error       bool
	FailedField string
	Tag         string
	Value       interface{}
}

func mapToErrorResponse(err validator.FieldError) errorResponse {
	return errorResponse{
		FailedField: err.Field(),
		Tag:         err.Tag(),
		Value:       err.Value(),
		Error:       true,
	}
}

func validateStruct(data interface{}) []errorResponse {
	validationErrors := []errorResponse{}

	if errs := validateSingleton.Struct(data); errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, mapToErrorResponse(err))
		}
	}

	return validationErrors
}

func getValidationMessage(validationErrors []errorResponse) string {
	if len(validationErrors) > 0 && validationErrors[0].Error {
		errMessages := strings.Builder{}
		for _, err := range validationErrors {
			str := fmt.Sprintf(
				"[%s]: '%v' | Needs to implement '%s'\n",
				err.FailedField,
				err.Value,
				err.Tag,
			)
			errMessages.WriteString(str)
		}
		return errMessages.String()
	}
	return ""
}

func ValidStruct(data interface{}) error {
	validationErrors := validateStruct(data)
	message := getValidationMessage(validationErrors)
	if message != "" {
		return errors.New(message)
	}
	return nil
}

func validateUUIDString(field interface{}) []errorResponse {
	validationErrors := []errorResponse{}

	if errs := validateSingleton.Var(field, "required,uuid4"); errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, mapToErrorResponse(err))
		}
	}

	return validationErrors
}

func ValidUUID(field interface{}) error {
	validationErrors := validateUUIDString(field)
	message := getValidationMessage(validationErrors)
	if message != "" {
		return errors.New(message)
	}
	return nil
}
