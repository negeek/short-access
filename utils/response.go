package utils

import (
	"encoding/json"
	"net/http"
)
type Response struct {
	Success  bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
// sends json response
func JsonResponse(w http.ResponseWriter, success bool, statusCode int, message string, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Success:  success,
		Message: message,
		Data:    data,
	})
}