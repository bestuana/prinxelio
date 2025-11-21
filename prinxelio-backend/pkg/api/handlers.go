package api

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type productDTO struct {
	ID                    int     `json:"id"`
	ProductName           string  `json:"product_name"`
	ProductImage          string  `json:"product_image"`
	ProductPrice          float64 `json:"product_price"`
	ProductDiscount       float64 `json:"product_discount"`
	ProductDiscountAmount int     `json:"product_discount_amount"`
	ProductDesc           string  `json:"product_desc"`
	ProductViewed         int     `json:"product_viewed"`
	ProductDownloaded     int     `json:"product_downloaded"`
	CategoryName          string  `json:"category_name"`
}

type categoryDTO struct {
	ID             int    `json:"id"`
	CategoryName   string `json:"category_name"`
	CategoryImages string `json:"category_images"`
	CategoryColor  string `json:"category_color"`
}

// format_response: Mengembalikan daftar produk untuk frontend (mendukung pencarian)
func (s *Server) handleProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteJSON(w, http.StatusMethodNotAllowed, APIResponse{Status: false, Message: "Metode tidak diizinkan", Data: nil})
		return
	}
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	var rows *sql.Rows
	var err error
	if q == "" {
		rows, err = s.DB.Query("SELECT p.id, p.product_name, p.product_image, p.product_price, p.product_discount, p.product_discount_amount, p.product_desc, p.product_viewed, p.product_downloaded, IFNULL(c.category_name,'') as category_name FROM product p LEFT JOIN category c ON c.id = p.product_category ORDER BY p.id DESC")
	} else {
		like := "%" + q + "%"
		rows, err = s.DB.Query("SELECT p.id, p.product_name, p.product_image, p.product_price, p.product_discount, p.product_discount_amount, p.product_desc, p.product_viewed, p.product_downloaded, IFNULL(c.category_name,'') as category_name FROM product p LEFT JOIN category c ON c.id = p.product_category WHERE p.product_name LIKE ? OR p.product_desc LIKE ? ORDER BY p.id DESC", like, like)
	}
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Status: false, Message: "Gagal memuat produk", Data: nil})
		return
	}
	defer rows.Close()
	products := []productDTO{}
	for rows.Next() {
		var p productDTO
		if err := rows.Scan(&p.ID, &p.ProductName, &p.ProductImage, &p.ProductPrice, &p.ProductDiscount, &p.ProductDiscountAmount, &p.ProductDesc, &p.ProductViewed, &p.ProductDownloaded, &p.CategoryName); err == nil {
			products = append(products, p)
		}
	}
	WriteJSON(w, http.StatusOK, APIResponse{Status: true, Message: "Berhasil memuat data", Data: products})
}

func (s *Server) handleCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteJSON(w, http.StatusMethodNotAllowed, APIResponse{Status: false, Message: "Metode tidak diizinkan", Data: nil})
		return
	}
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	var rows *sql.Rows
	var err error
	if q == "" {
		rows, err = s.DB.Query("SELECT id, category_name, IFNULL(category_images,''), IFNULL(category_color,'') FROM category ORDER BY id ASC")
	} else {
		like := "%" + q + "%"
		rows, err = s.DB.Query("SELECT id, category_name, IFNULL(category_images,''), IFNULL(category_color,'') FROM category WHERE category_name LIKE ? ORDER BY id ASC", like)
	}
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Status: false, Message: "Gagal memuat kategori", Data: nil})
		return
	}
	defer rows.Close()
	list := []categoryDTO{}
	for rows.Next() {
		var c categoryDTO
		if err := rows.Scan(&c.ID, &c.CategoryName, &c.CategoryImages, &c.CategoryColor); err == nil {
			list = append(list, c)
		}
	}
	WriteJSON(w, http.StatusOK, APIResponse{Status: true, Message: "Berhasil", Data: list})
}

type sendOTPReq struct {
	Phone string `json:"phone"`
}

// format_send: Kirim envelope SEND_OTP ke n8n; format_response untuk frontend
func (s *Server) handleSendOTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteJSON(w, http.StatusMethodNotAllowed, APIResponse{Status: false, Message: "Metode tidak diizinkan", Data: nil})
		return
	}
	var req sendOTPReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Phone == "" {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Status: false, Message: "Nomor WhatsApp wajib diisi", Data: nil})
		return
	}
	otp := generateOTP()
	now := time.Now()
	_, err := s.DB.Exec("INSERT INTO users (phone_number, otp_code, otp_created_at, status) VALUES (?, ?, ?, 1) ON DUPLICATE KEY UPDATE otp_code=VALUES(otp_code), otp_created_at=VALUES(otp_created_at)", req.Phone, otp, now)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Status: false, Message: "Gagal menyimpan OTP", Data: nil})
		return
	}
	signature := signOTP(req.Phone, otp, now.Unix())
	payload := map[string]interface{}{
		"phone":     req.Phone,
		"otp":       otp,
		"signature": signature,
	}
	log.Printf("otp_send payload=%v", payload)
	resp, err := SendToWebhook(os.Getenv("N8N_WEBHOOK_URL_OTP"), "SEND_OTP", payload, N8NBearer())
	if err != nil {
		log.Printf("otp_send webhook_error err=%v", err)
		WriteJSON(w, http.StatusOK, APIResponse{Status: true, Message: "OTP diproses, webhook tidak tersedia", Data: map[string]interface{}{"expires_in": os.Getenv("TIME_OTP")}})
		return
	}
	body, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	log.Printf("otp_send webhook_status=%d body=%s", resp.StatusCode, string(body))
	var hook struct {
		Status bool `json:"status"`
		Data   struct {
			Delivered bool   `json:"delivered"`
			Message   string `json:"message"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &hook); err != nil {
		log.Printf("otp_send unmarshal_error err=%v", err)
		WriteJSON(w, http.StatusOK, APIResponse{Status: false, Message: "Respons webhook tidak valid", Data: map[string]interface{}{"webhook_raw": string(body), "webhook_status": resp.StatusCode}})
		return
	}
	if !hook.Status || !hook.Data.Delivered {
		WriteJSON(w, http.StatusOK, APIResponse{Status: false, Message: hook.Data.Message, Data: map[string]interface{}{"expires_in": os.Getenv("TIME_OTP"), "delivered": hook.Data.Delivered, "message": hook.Data.Message, "webhook_status": resp.StatusCode}})
		return
	}
	WriteJSON(w, http.StatusOK, APIResponse{Status: true, Message: "OTP dikirim", Data: map[string]interface{}{"expires_in": os.Getenv("TIME_OTP"), "delivered": hook.Data.Delivered, "message": hook.Data.Message, "webhook_status": resp.StatusCode}})
}

type verifyOTPReq struct {
	Phone string `json:"phone"`
	OTP   string `json:"otp"`
}

// format_response: Verifikasi OTP dan kembalikan JWT untuk frontend
func (s *Server) handleVerifyOTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteJSON(w, http.StatusMethodNotAllowed, APIResponse{Status: false, Message: "Metode tidak diizinkan", Data: nil})
		return
	}
	var req verifyOTPReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Phone == "" || req.OTP == "" {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Status: false, Message: "Phone dan OTP wajib diisi", Data: nil})
		return
	}
	var code sql.NullString
	var createdAt sql.NullTime
	err := s.DB.QueryRow("SELECT otp_code, otp_created_at FROM users WHERE phone_number = ?", req.Phone).Scan(&code, &createdAt)
	if err != nil || !code.Valid || !createdAt.Valid {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Status: false, Message: "OTP tidak ditemukan", Data: nil})
		return
	}
	if code.String != req.OTP {
		WriteJSON(w, http.StatusUnauthorized, APIResponse{Status: false, Message: "OTP salah", Data: nil})
		return
	}
	ttl, _ := strconv.Atoi(os.Getenv("TIME_OTP"))
	if time.Since(createdAt.Time) > time.Duration(ttl)*time.Second {
		WriteJSON(w, http.StatusUnauthorized, APIResponse{Status: false, Message: "OTP kedaluwarsa", Data: nil})
		return
	}
	token, err := s.issueJWT(req.Phone)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Status: false, Message: "Gagal membuat token", Data: nil})
		return
	}
	_, _ = s.DB.Exec("UPDATE users SET last_login = NOW(), otp_code = NULL WHERE phone_number = ?", req.Phone)
	WriteJSON(w, http.StatusOK, APIResponse{Status: true, Message: "Verifikasi berhasil", Data: map[string]string{"token": token}})
}

type createTxReq struct {
	ProductID int    `json:"product_id"`
	Phone     string `json:"phone"`
}

// format_send: Kirim envelope CREATE_TRANSACTION ke n8n; format_response untuk QR
func (s *Server) handleCreateTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteJSON(w, http.StatusMethodNotAllowed, APIResponse{Status: false, Message: "Metode tidak diizinkan", Data: nil})
		return
	}
	if !s.validateJWT(r.Header.Get("Authorization")) {
		WriteJSON(w, http.StatusUnauthorized, APIResponse{Status: false, Message: "Token tidak valid", Data: nil})
		return
	}
	var req createTxReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ProductID == 0 || req.Phone == "" {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Status: false, Message: "Data tidak lengkap", Data: nil})
		return
	}
	var base, disc float64
	var discountAmt int
	var name string
	err := s.DB.QueryRow("SELECT product_price, product_discount, product_discount_amount, product_name FROM product WHERE id = ?", req.ProductID).Scan(&base, &disc, &discountAmt, &name)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Status: false, Message: "Produk tidak ditemukan", Data: nil})
		return
	}
	amount := base
	if disc > 0 {
		amount = disc
	}
	if disc == 0 || discountAmt == 100 || base == 0 {
		amount = 0
	}
	log.Printf("tx_create product_id=%d base=%.2f disc=%.2f amount=%.2f", req.ProductID, base, disc, amount)
	// Kirim ke webhook tanpa reference; tunggu respon dan baru catat
	payload := map[string]interface{}{
		"amount":         amount,
		"product_name":   name,
		"customer_phone": req.Phone,
	}
	log.Printf("tx_create webhook_payload=%v", payload)
	resp, err := SendToWebhook(os.Getenv("N8N_WEBHOOK_URL_TRANSACTION_CREATE"), "CREATE_TRANSACTION", payload, N8NBearer())
	if err != nil {
		log.Printf("tx_create webhook_error err=%v", err)
		WriteJSON(w, http.StatusOK, APIResponse{Status: false, Message: "Webhook transaksi tidak tersedia", Data: nil})
		return
	}
	body, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	log.Printf("tx_create webhook_status=%d body=%s", resp.StatusCode, string(body))
	var hook struct {
		Success bool `json:"success"`
		Data    struct {
			Reference     string  `json:"reference"`
			QrURL         string  `json:"qr_url"`
			MerchantRef   string  `json:"merchant_ref"`
			CustomerName  string  `json:"customer_name"`
			CustomerEmail string  `json:"customer_email"`
			CustomerPhone string  `json:"customer_phone"`
			Amount        float64 `json:"amount"`
			TotalFee      float64 `json:"total_fee"`
			Status        string  `json:"status"`
			ExpiredTime   int64   `json:"expired_time"`
			Message       string  `json:"message"`
		} `json:"data"`
	}
	decErr := json.Unmarshal(body, &hook)
	if decErr != nil || !hook.Success {
		msg := "Gagal Membuat transaksi!"
		// jika ada message dalam data, gunakan
		if decErr == nil && hook.Data.Message != "" {
			msg = hook.Data.Message
		}
		if decErr != nil {
			log.Printf("tx_create unmarshal_error err=%v", decErr)
		}
		WriteJSON(w, http.StatusOK, APIResponse{Status: false, Message: msg, Data: map[string]interface{}{"webhook_raw": string(body), "webhook_status": resp.StatusCode}})
		return
	}
	// transaksi gratis (amount 0) diproses langsung
	if amount <= 0 || strings.EqualFold(hook.Data.Message, "Transaksi Gratis!") || (hook.Data.Amount == 0 && hook.Data.Reference == "") {
		ensureTransactionStatusEnum(s.DB)
		_, _ = s.DB.Exec("INSERT INTO users (phone_number, status) VALUES (?, 1) ON DUPLICATE KEY UPDATE status=VALUES(status)", req.Phone)
		var userID int
		_ = s.DB.QueryRow("SELECT id FROM users WHERE phone_number = ?", req.Phone).Scan(&userID)
		b := make([]byte, 5)
		_, _ = rand.Read(b)
		freeRef := "FREE-" + hex.EncodeToString(b)
		_, err = s.DB.Exec(
			"INSERT INTO transactions (user_id, product_id, base_price, admin_fee, total_amount, status, merchant_ref, expired_time, reference, link_qr) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			userID, req.ProductID, base, 0, 0, "PAID", "", 0, freeRef, "",
		)
		if err != nil {
			ensureTransactionStatusEnum(s.DB)
			_, _ = s.DB.Exec(
				"INSERT INTO transactions (user_id, product_id, base_price, admin_fee, total_amount, status, merchant_ref, expired_time, reference, link_qr) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
				userID, req.ProductID, base, 0, 0, "PAID", "", 0, freeRef, "",
			)
		}
		_, _ = s.DB.Exec("UPDATE product SET product_downloaded = product_downloaded + 1 WHERE id = ?", req.ProductID)
		WriteJSON(w, http.StatusOK, APIResponse{Status: true, Message: "Transaksi Gratis!", Data: map[string]interface{}{"message": "Transaksi Gratis!"}})
		return
	}
	// Pastikan enum status sesuai dan user ada
	ensureTransactionStatusEnum(s.DB)
	_, _ = s.DB.Exec("INSERT INTO users (phone_number, status) VALUES (?, 1) ON DUPLICATE KEY UPDATE status=VALUES(status)", req.Phone)
	var userID int
	_ = s.DB.QueryRow("SELECT id FROM users WHERE phone_number = ?", req.Phone).Scan(&userID)
	// Sanitasi qr_url (hilangkan backticks dan spasi)
	qrURL := strings.Trim(hook.Data.QrURL, " `")
	status := normalizeStatusForDB(hook.Data.Status)
	// Catat ke database setelah sukses
	_, err = s.DB.Exec(
		"INSERT INTO transactions (user_id, product_id, base_price, admin_fee, total_amount, status, merchant_ref, expired_time, reference, link_qr) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		userID, req.ProductID, base, hook.Data.TotalFee, hook.Data.Amount, status, hook.Data.MerchantRef, hook.Data.ExpiredTime, hook.Data.Reference, qrURL,
	)
	if err != nil {
		log.Printf("tx_create insert_error err=%v phone=%s product_id=%d base=%.2f fee=%.2f amount=%.2f status=%s merchant_ref=%s expired=%d reference=%s qr=%s",
			err, req.Phone, req.ProductID, base, hook.Data.TotalFee, hook.Data.Amount, status, hook.Data.MerchantRef, hook.Data.ExpiredTime, hook.Data.Reference, qrURL)
		// Coba alter enum, lalu sisipkan ulang sekali
		ensureTransactionStatusEnum(s.DB)
		_, retryErr := s.DB.Exec(
			"INSERT INTO transactions (user_id, product_id, base_price, admin_fee, total_amount, status, merchant_ref, expired_time, reference, link_qr) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			userID, req.ProductID, base, hook.Data.TotalFee, hook.Data.Amount, status, hook.Data.MerchantRef, hook.Data.ExpiredTime, hook.Data.Reference, qrURL,
		)
		if retryErr != nil {
			log.Printf("tx_create insert_retry_error err=%v", retryErr)
			WriteJSON(w, http.StatusInternalServerError, APIResponse{Status: false, Message: "Gagal menyimpan transaksi", Data: nil})
			return
		}
	}
	log.Printf("tx_create insert_ok reference=%s", hook.Data.Reference)
	WriteJSON(w, http.StatusOK, APIResponse{Status: true, Message: "Transaksi dibuat", Data: map[string]interface{}{
		"reference":    hook.Data.Reference,
		"qr_url":       qrURL,
		"merchant_ref": hook.Data.MerchantRef,
		"status":       hook.Data.Status,
		"expired_time": hook.Data.ExpiredTime,
		"amount":       hook.Data.Amount,
		"total_fee":    hook.Data.TotalFee,
	}})
}

// format_response: Mengembalikan status transaksi untuk tombol cek status
func (s *Server) handleGetTransactionStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteJSON(w, http.StatusMethodNotAllowed, APIResponse{Status: false, Message: "Metode tidak diizinkan", Data: nil})
		return
	}
	phone, ok := s.subjectFromJWT(r.Header.Get("Authorization"))
	if !ok {
		WriteJSON(w, http.StatusUnauthorized, APIResponse{Status: false, Message: "Token tidak valid", Data: nil})
		return
	}
	ref := r.URL.Query().Get("reference")
	if ref == "" {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Status: false, Message: "Reference wajib diisi", Data: nil})
		return
	}
	var status string
	err := s.DB.QueryRow("SELECT t.status FROM transactions t JOIN users u ON t.user_id = u.id WHERE t.reference = ? AND u.phone_number = ?", ref, phone).Scan(&status)
	if err != nil {
		WriteJSON(w, http.StatusNotFound, APIResponse{Status: false, Message: "Transaksi tidak ditemukan", Data: nil})
		return
	}
	WriteJSON(w, http.StatusOK, APIResponse{Status: true, Message: "Berhasil", Data: map[string]string{"status": status}})
}

// format_response: Batalkan transaksi dan siarkan status via WebSocket
func (s *Server) handleCancelTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteJSON(w, http.StatusMethodNotAllowed, APIResponse{Status: false, Message: "Metode tidak diizinkan", Data: nil})
		return
	}
	phone, ok := s.subjectFromJWT(r.Header.Get("Authorization"))
	if !ok {
		WriteJSON(w, http.StatusUnauthorized, APIResponse{Status: false, Message: "Token tidak valid", Data: nil})
		return
	}
	ref := r.URL.Query().Get("reference")
	if ref == "" {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Status: false, Message: "Reference wajib diisi", Data: nil})
		return
	}
	_, err := s.DB.Exec("UPDATE transactions t JOIN users u ON t.user_id = u.id SET t.status = 'FAILED' WHERE t.reference = ? AND u.phone_number = ? AND t.status = 'UNPAID'", ref, phone)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Status: false, Message: "Gagal membatalkan", Data: nil})
		return
	}
	s.Hub.Broadcast(ref, map[string]string{"status": "FAILED"})
	WriteJSON(w, http.StatusOK, APIResponse{Status: true, Message: "Dibatalkan", Data: map[string]string{"reference": ref}})
}

// format_response: Endpoint webhook n8n untuk update status (PAID/FAILED/EXPIRED)
func (s *Server) handleWebhookPaymentStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteJSON(w, http.StatusMethodNotAllowed, APIResponse{Status: false, Message: "Metode tidak diizinkan", Data: nil})
		return
	}
	auth := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if auth == "" || auth != os.Getenv("N8N_SECRETKEY_JWT") {
		WriteJSON(w, http.StatusUnauthorized, APIResponse{Status: false, Message: "Tidak diizinkan", Data: nil})
		return
	}
	var body struct {
		Reference string `json:"reference"`
		Status    string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Reference == "" {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Status: false, Message: "Data tidak valid", Data: nil})
		return
	}
	switch body.Status {
	case "PAID":
		_, _ = s.DB.Exec("UPDATE transactions SET status='PAID' WHERE reference = ?", body.Reference)
		_, _ = s.DB.Exec("UPDATE product p JOIN transactions t ON p.id = t.product_id SET p.product_downloaded = p.product_downloaded + 1 WHERE t.reference = ?", body.Reference)
	case "FAILED", "EXPIRED":
		_, _ = s.DB.Exec("UPDATE transactions SET status=? WHERE reference = ?", body.Status, body.Reference)
	case "PENDING":
		_, _ = s.DB.Exec("UPDATE transactions SET status='UNPAID' WHERE reference = ?", body.Reference)
	default:
	}
	s.Hub.Broadcast(body.Reference, map[string]string{"status": body.Status})
	WriteJSON(w, http.StatusOK, APIResponse{Status: true, Message: "OK", Data: nil})
}

func (s *Server) basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != os.Getenv("AUTH_USER") || pass != os.Getenv("AUTH_PASS") {
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			WriteJSON(w, http.StatusUnauthorized, APIResponse{Status: false, Message: "Unauthorized", Data: nil})
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (s *Server) handleAdminDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteJSON(w, http.StatusMethodNotAllowed, APIResponse{Status: false, Message: "Metode tidak diizinkan", Data: nil})
		return
	}
	var views, purchases, paid, failed int
	_ = s.DB.QueryRow("SELECT IFNULL(SUM(product_viewed),0) FROM product").Scan(&views)
	_ = s.DB.QueryRow("SELECT COUNT(*) FROM transactions").Scan(&purchases)
	_ = s.DB.QueryRow("SELECT COUNT(*) FROM transactions WHERE status='PAID'").Scan(&paid)
	_ = s.DB.QueryRow("SELECT COUNT(*) FROM transactions WHERE status IN ('FAILED','EXPIRED')").Scan(&failed)
	WriteJSON(w, http.StatusOK, APIResponse{Status: true, Message: "Berhasil", Data: map[string]int{"views": views, "purchases": purchases, "paid": paid, "failed": failed}})
}

func (s *Server) handleAdminProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleProducts(w, r)
	case http.MethodPost:
		var p productDTO
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil || p.ProductName == "" || p.ProductPrice == 0 {
			WriteJSON(w, http.StatusBadRequest, APIResponse{Status: false, Message: "Data produk tidak valid", Data: nil})
			return
		}
		_, err := s.DB.Exec("INSERT INTO product (product_name, product_image, product_price, product_discount, product_discount_amount, product_desc) VALUES (?, ?, ?, ?, ?, ?)", p.ProductName, p.ProductImage, p.ProductPrice, p.ProductDiscount, p.ProductDiscountAmount, p.ProductDesc)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, APIResponse{Status: false, Message: "Gagal menambah produk", Data: nil})
			return
		}
		WriteJSON(w, http.StatusOK, APIResponse{Status: true, Message: "Produk ditambah", Data: nil})
	case http.MethodPut:
		var p productDTO
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil || p.ID == 0 {
			WriteJSON(w, http.StatusBadRequest, APIResponse{Status: false, Message: "ID produk wajib", Data: nil})
			return
		}
		_, err := s.DB.Exec("UPDATE product SET product_name=?, product_image=?, product_price=?, product_discount=?, product_discount_amount=?, product_desc=? WHERE id=?", p.ProductName, p.ProductImage, p.ProductPrice, p.ProductDiscount, p.ProductDiscountAmount, p.ProductDesc, p.ID)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, APIResponse{Status: false, Message: "Gagal mengedit produk", Data: nil})
			return
		}
		WriteJSON(w, http.StatusOK, APIResponse{Status: true, Message: "Produk diedit", Data: nil})
	case http.MethodDelete:
		idStr := r.URL.Query().Get("id")
		id, _ := strconv.Atoi(idStr)
		if id == 0 {
			WriteJSON(w, http.StatusBadRequest, APIResponse{Status: false, Message: "ID produk wajib", Data: nil})
			return
		}
		_, err := s.DB.Exec("DELETE FROM product WHERE id = ?", id)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, APIResponse{Status: false, Message: "Gagal menghapus produk", Data: nil})
			return
		}
		WriteJSON(w, http.StatusOK, APIResponse{Status: true, Message: "Produk dihapus", Data: nil})
	default:
		WriteJSON(w, http.StatusMethodNotAllowed, APIResponse{Status: false, Message: "Metode tidak diizinkan", Data: nil})
	}
}

// format_response: Mencatat view produk
func (s *Server) handleProductView(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteJSON(w, http.StatusMethodNotAllowed, APIResponse{Status: false, Message: "Metode tidak diizinkan", Data: nil})
		return
	}
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	if id == 0 {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Status: false, Message: "ID produk wajib", Data: nil})
		return
	}
	_, err := s.DB.Exec("UPDATE product SET product_viewed = product_viewed + 1 WHERE id = ?", id)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Status: false, Message: "Gagal update view", Data: nil})
		return
	}
	WriteJSON(w, http.StatusOK, APIResponse{Status: true, Message: "OK", Data: nil})
}

func (s *Server) apiSecretToday() string {
	key := os.Getenv("PUBLIC_API_SECRET_KEY")
	day := time.Now().Format("2006-01-02")
	sum := sha256.Sum256([]byte(key + "|" + day))
	return hex.EncodeToString(sum[:])
}

func (s *Server) apiSecretGuard(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("sk")
		if err != nil || c.Value == "" || c.Value != s.apiSecretToday() {
			WriteJSON(w, http.StatusUnauthorized, APIResponse{Status: false, Message: "Tidak diizinkan", Data: nil})
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (s *Server) issueJWT(subject string) (string, error) {
	dur, _ := strconv.Atoi(os.Getenv("TIME_JWT"))
	claims := jwt.MapClaims{
		"sub": subject,
		"exp": time.Now().Add(time.Duration(dur) * time.Second).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("PUBLIC_API_SECRET_KEY")))
}

func (s *Server) validateJWT(authHeader string) bool {
	parts := strings.Split(strings.TrimSpace(authHeader), " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return false
	}
	tokenStr := parts[1]
	_, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid")
		}
		return []byte(os.Getenv("PUBLIC_API_SECRET_KEY")), nil
	})
	return err == nil
}

func (s *Server) subjectFromJWT(authHeader string) (string, bool) {
	parts := strings.Split(strings.TrimSpace(authHeader), " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", false
	}
	tokenStr := parts[1]
	t, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid")
		}
		return []byte(os.Getenv("PUBLIC_API_SECRET_KEY")), nil
	})
	if err != nil || !t.Valid {
		return "", false
	}
	if claims, ok := t.Claims.(jwt.MapClaims); ok {
		if sub, ok := claims["sub"].(string); ok {
			return sub, true
		}
	}
	return "", false
}

func generateOTP() string {
	b := make([]byte, 3)
	_, _ = rand.Read(b)
	n := int(b[0])<<16 | int(b[1])<<8 | int(b[2])
	return fmt.Sprintf("%06d", n%1000000)
}

func signOTP(phone, otp string, ts int64) string {
	key := []byte(os.Getenv("OTP_SECRET_KEY"))
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(fmt.Sprintf("%s|%s|%d", phone, otp, ts)))
	return hex.EncodeToString(mac.Sum(nil))
}

// Memastikan kolom enum status pada tabel transactions sudah: UNPAID, PAID, FAILED, EXPIRED
func ensureTransactionStatusEnum(db *sql.DB) {
	var colType sql.NullString
	err := db.QueryRow(
		"SELECT COLUMN_TYPE FROM information_schema.columns WHERE table_schema = ? AND table_name = 'transactions' AND column_name = 'status'",
		os.Getenv("DB_NAME"),
	).Scan(&colType)
	if err != nil {
		log.Printf("ensure_enum query_error err=%v", err)
		return
	}
	if !colType.Valid || !strings.Contains(colType.String, "'UNPAID'") || !strings.Contains(colType.String, "'PAID'") || !strings.Contains(colType.String, "'FAILED'") || !strings.Contains(colType.String, "'EXPIRED'") {
		_, alterErr := db.Exec("ALTER TABLE transactions MODIFY COLUMN status ENUM('UNPAID','PAID','FAILED','EXPIRED') DEFAULT 'UNPAID'")
		if alterErr != nil {
			log.Printf("ensure_enum alter_error err=%v", alterErr)
		} else {
			log.Printf("ensure_enum altered to UNPAID,PAID,FAILED,EXPIRED")
		}
	}
}

// Menormalkan status ke salah satu dari UNPAID, PAID, FAILED, EXPIRED
func normalizeStatusForDB(s string) string {
	x := strings.ToUpper(strings.TrimSpace(s))
	switch x {
	case "UNPAID", "PAID", "FAILED", "EXPIRED":
		return x
	case "PENDING":
		return "UNPAID"
	default:
		return "UNPAID"
	}
}
