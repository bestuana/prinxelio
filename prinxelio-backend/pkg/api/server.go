package api

import (
	"database/sql"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Server struct {
	DB  *sql.DB
	Mux *http.ServeMux
	Hub *TxHub
}

func NewServer(db *sql.DB) *Server {
	s := &Server{DB: db, Mux: http.NewServeMux(), Hub: NewTxHub()}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.Mux.HandleFunc("/api/products", s.apiSecretGuard(s.handleProducts))
	s.Mux.HandleFunc("/api/categories", s.apiSecretGuard(s.handleCategories))
	s.Mux.HandleFunc("/api/otp/send", s.apiSecretGuard(s.handleSendOTP))
	s.Mux.HandleFunc("/api/otp/verify", s.apiSecretGuard(s.handleVerifyOTP))
	s.Mux.HandleFunc("/api/transactions/create", s.handleCreateTransaction)
	s.Mux.HandleFunc("/api/transactions/status", s.handleGetTransactionStatus)
	s.Mux.HandleFunc("/api/transactions/cancel", s.handleCancelTransaction)
	s.Mux.HandleFunc("/api/ws/transactions", s.handleWsTransactions)
	s.Mux.HandleFunc("/api/webhook/payment-status", s.handleWebhookPaymentStatus)
	s.Mux.HandleFunc("/api/products/view", s.apiSecretGuard(s.handleProductView))

	s.Mux.HandleFunc("/admin/dashboard", s.basicAuth(s.handleAdminDashboard))
	s.Mux.HandleFunc("/admin/products", s.basicAuth(s.handleAdminProducts))

	s.Mux.Handle("/", s.staticWithCookie(http.FileServer(http.Dir("public"))))
	s.Mux.HandleFunc("/admin", s.adminProxy)
	s.Mux.HandleFunc("/admin/", s.adminProxy)
	s.Mux.HandleFunc("/admin.php", s.adminProxy)
}

func (s *Server) staticWithCookie(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "sk", Value: s.apiSecretToday(), Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode})
		p := r.URL.Path
		if p == "/" || p == "" {
			http.ServeFile(w, r, filepath.Join("public", "index.html"))
			return
		}
		full := filepath.Join("public", filepath.Clean(p))
		if filepath.Ext(full) == ".php" {
			http.NotFound(w, r)
			return
		}
		if fi, err := os.Stat(full); err == nil {
			if fi.IsDir() {
				h.ServeHTTP(w, r)
				return
			}
			h.ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join("public", "index.html"))
	})
}

func (s *Server) adminProxy(w http.ResponseWriter, r *http.Request) {
    target := os.Getenv("ADMIN_PHP_URL")
    if target == "" { target = "http://admin" }
    u, err := url.Parse(target)
    if err != nil { http.NotFound(w, r); return }
    origHost := r.Host
    scheme := "http"
    if h := r.Header.Get("X-Forwarded-Proto"); h != "" { scheme = h }
    proxy := httputil.NewSingleHostReverseProxy(u)
    proxy.Director = func(req *http.Request) {
        req.URL.Scheme = u.Scheme
        req.URL.Host = u.Host
        p := req.URL.Path
        if p == "/admin" || p == "/admin/" || p == "/admin.php" || p == "/admin.php/" {
            req.URL.Path = "/admin.php"
        } else if strings.HasPrefix(p, "/admin.php") {
            q := strings.TrimPrefix(p, "/admin.php")
            if q == "" { q = "/admin.php" }
            if !strings.HasPrefix(q, "/") { q = "/" + q }
            req.URL.Path = q
        } else if strings.HasPrefix(p, "/admin/") {
            q := strings.TrimPrefix(p, "/admin")
            if q == "" { q = "/" }
            if !strings.HasPrefix(q, "/") { q = "/" + q }
            req.URL.Path = q
        } else {
            req.URL.Path = "/admin.php"
        }
        req.Host = u.Host
        req.Header.Set("X-Forwarded-Host", origHost)
        req.Header.Set("X-Forwarded-Proto", scheme)
    }
    proxy.ModifyResponse = func(resp *http.Response) error {
        loc := resp.Header.Get("Location")
        if loc == "" { return nil }
        lu, err := url.Parse(loc)
        if err != nil { return nil }
        if lu.Host == u.Host || lu.Host == "" {
            np := lu.Path
            if !strings.HasPrefix(np, "/admin") {
                np = "/admin" + np
            }
            lu.Scheme = scheme
            lu.Host = origHost
            lu.Path = np
            resp.Header.Set("Location", lu.String())
        }
        return nil
    }
    proxy.ServeHTTP(w, r)
}
