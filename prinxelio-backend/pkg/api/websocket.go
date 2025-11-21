package api

import (
    "encoding/json"
    "net/http"
    "sync"

    "github.com/gorilla/websocket"
)

type TxHub struct {
    mu    sync.RWMutex
    conns map[string]map[*websocket.Conn]struct{}
    up    websocket.Upgrader
}

func NewTxHub() *TxHub {
    return &TxHub{conns: make(map[string]map[*websocket.Conn]struct{}), up: websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}}
}

func (s *Server) handleWsTransactions(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet { WriteJSON(w, http.StatusMethodNotAllowed, APIResponse{Status:false, Message:"Metode tidak diizinkan", Data:nil}); return }
    ref := r.URL.Query().Get("reference")
    if ref == "" { WriteJSON(w, http.StatusBadRequest, APIResponse{Status:false, Message:"Reference wajib diisi", Data:nil}); return }
    conn, err := s.Hub.up.Upgrade(w, r, nil)
    if err != nil { return }
    s.Hub.add(ref, conn)
    defer s.Hub.remove(ref, conn)
    s.Hub.send(conn, map[string]string{"status":"CONNECTED"})
    for {
        _, _, err := conn.ReadMessage()
        if err != nil { break }
    }
}

func (h *TxHub) add(ref string, c *websocket.Conn) {
    h.mu.Lock()
    defer h.mu.Unlock()
    if _, ok := h.conns[ref]; !ok { h.conns[ref] = make(map[*websocket.Conn]struct{}) }
    h.conns[ref][c] = struct{}{}
}

func (h *TxHub) remove(ref string, c *websocket.Conn) {
    h.mu.Lock()
    defer h.mu.Unlock()
    if cs, ok := h.conns[ref]; ok { delete(cs, c) }
    _ = c.Close()
}

func (h *TxHub) Broadcast(ref string, v interface{}) {
    h.mu.RLock()
    cs := h.conns[ref]
    h.mu.RUnlock()
    if cs == nil { return }
    for c := range cs { h.send(c, v) }
}

func (h *TxHub) send(c *websocket.Conn, v interface{}) {
    b, _ := json.Marshal(v)
    _ = c.WriteMessage(websocket.TextMessage, b)
}
