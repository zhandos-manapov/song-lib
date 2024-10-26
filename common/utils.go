package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func MakeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle the error
			var e *apiError
			if errors.As(err, &e) {
				WriteJSON(w, e.Code, newApiError(e.Code, e.Error(), e.Unwrap()))
			} else {
				WriteJSON(w, http.StatusInternalServerError, NewInternalServerError(err.Error(), nil))
			}
		}
	}
}

func ParseJSON(r io.Reader, payload any) error {
	return json.NewDecoder(r).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func ParseBody[T interface{}](r *http.Request, payload *T) (*T, error) {
	if r.Body == nil {
		return nil, fmt.Errorf("missing request body")
	}

	if err := ParseJSON(r.Body, payload); err != nil {
		return nil, err
	}

	if err := ValidStruct(*payload); err != nil {
		return nil, err
	}

	return payload, nil
}
