package common

import (
	"errors"
	"fmt"
	"net/http"
)

type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	err     error
	Details string `json:"details"`
}

func newApiError(code int, message string, err error) *apiError {
	apiError := &apiError{
		Code:    code,
		Message: message,
	}
	if err == nil {
		apiError.err = errors.New(message)
		apiError.Details = apiError.err.Error()
	} else {
		apiError.err = err
		apiError.Details = err.Error()
	}

	return apiError
}

func (e *apiError) Error() string {
	return fmt.Sprintf("%d - %s", e.Code, e.Message)
}

func (e *apiError) Unwrap() error {
	return e.err
}

func NewBadRequestError(message string, err error) *apiError {
	return newApiError(http.StatusBadRequest, message, err)
}

func NewNotFoundError(message string, err error) *apiError {
	return newApiError(http.StatusNotFound, message, err)
}

func NewInternalServerError(message string, err error) *apiError {
	return newApiError(http.StatusInternalServerError, message, err)
}
