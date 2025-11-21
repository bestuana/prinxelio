package api

import (
    "encoding/json"
    "net/http"
)

// format_response: Semua handler API harus mengembalikan JSON sesuai struktur APIResponse
type APIResponse struct {
    Status  bool        `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

func WriteJSON(w http.ResponseWriter, statusCode int, payload APIResponse) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    _ = json.NewEncoder(w).Encode(payload)
}
