package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type PaginatedResponse struct {
	Data  interface{} `json:"data"`
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
}

var validate = validator.New()

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Message: message})
}

func getIntParam(r *http.Request, param string) (int, error) {
	paramStr := chi.URLParam(r, param)
	return strconv.Atoi(paramStr)
}

func getQueryInt(r *http.Request, param string, defaultValue int) int {
	paramStr := r.URL.Query().Get(param)
	if paramStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(paramStr)
	if err != nil {
		return defaultValue
	}

	return value
}
