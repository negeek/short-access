package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/negeek/short-access/apperr"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// JsonResponse sends a JSON response with the given status code.
func JsonResponse(w http.ResponseWriter, success bool, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Success: success,
		Message: message,
		Data:    data,
	})
}

// RespondError turns a service error into the right status code and a
// client-safe message. Anything that is not an *apperr.Error is treated as an
// unexpected failure.
func RespondError(w http.ResponseWriter, err error) {
	var appErr *apperr.Error
	if errors.As(err, &appErr) {
		JsonResponse(w, false, statusFor(appErr.Kind), appErr.Message, nil)
		return
	}
	JsonResponse(w, false, http.StatusInternalServerError, "Something went wrong. Try again.", nil)
}

func statusFor(kind apperr.Kind) int {
	switch kind {
	case apperr.KindBadRequest:
		return http.StatusBadRequest
	case apperr.KindUnauthorized:
		return http.StatusUnauthorized
	case apperr.KindNotFound:
		return http.StatusNotFound
	case apperr.KindConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
