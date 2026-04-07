package handler

import (
	"encoding/json"
	"net/http"
)

type UserParams struct {
	Email    string `json:"name"`
	Password string `json:"password"`
}

type ClientParams struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	InvoiceLink string `json:"invoice_link"`
	SendDate    string `json:"send_date"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}
