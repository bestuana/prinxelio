package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type WebhookEnvelope struct {
	Source    string      `json:"source"`
	Timestamp int64       `json:"timestamp"`
	Action    string      `json:"action"`
	Payload   interface{} `json:"payload"`
}

// format_send: Gunakan fungsi ini untuk mengirim data ke webhook n8n dengan envelope yang konsisten
func SendToWebhook(url string, action string, payload interface{}, authBearer string) (*http.Response, error) {
	env := WebhookEnvelope{
		Source:    "prinxelio_web",
		Timestamp: time.Now().Unix(),
		Action:    action,
		Payload:   payload,
	}
	body, _ := json.Marshal(env)
	log.Printf("webhook_send url=%s action=%s body=%s", url, action, string(body))
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Prinxelio-Backend/1.0")
	if authBearer != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authBearer))
	}
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	log.Printf("webhook_resp url=%s status=%d", url, resp.StatusCode)
	return resp, nil
}

func N8NBearer() string {
	secret := os.Getenv("N8N_SECRETKEY_JWT")
	claims := jwt.MapClaims{
		"iss": "prinxelio_web",
		"sub": "webhook-call",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(2 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(secret))
	return signed
}
