package presence

import (
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Viewer struct {
	UserID        string `json:"userID"`
	InformantName string `json:"informantName"`
	CustomerID    int64  `json:"customerID"`
}

type Message struct {
	Type string `json:"type"`
	Viewer
}

type Hub struct {
	clients map[*websocket.Conn]string
	viewers map[string]Viewer
	mu      sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients: map[*websocket.Conn]string{},
		viewers: map[string]Viewer{},
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Hub) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer h.disconnect(conn)

	h.connect(conn)

	for {
		var message Message

		err = conn.ReadJSON(&message)
		if err != nil {
			return
		}

		message.Type = strings.TrimSpace(message.Type)
		message.UserID = strings.TrimSpace(message.UserID)
		message.InformantName = strings.TrimSpace(message.InformantName)

		switch message.Type {
		case "", "view":
			if message.UserID == "" || message.InformantName == "" || message.CustomerID <= 0 {
				continue
			}

			h.update(conn, message.Viewer)
		case "leave":
			h.leave(conn, message.UserID)
		default:
			continue
		}
	}
}

func (h *Hub) connect(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[conn] = ""
	h.broadcastLocked()
}

func (h *Hub) update(conn *websocket.Conn, viewer Viewer) {
	h.mu.Lock()
	defer h.mu.Unlock()

	previousUserID := h.clients[conn]
	if previousUserID != "" && previousUserID != viewer.UserID {
		delete(h.viewers, previousUserID)
	}

	h.clients[conn] = viewer.UserID
	h.viewers[viewer.UserID] = viewer
	h.broadcastLocked()
}

func (h *Hub) leave(conn *websocket.Conn, userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if userID == "" {
		userID = h.clients[conn]
	}

	if userID == "" {
		return
	}

	for client, clientUserID := range h.clients {
		if clientUserID == userID {
			h.clients[client] = ""
		}
	}

	delete(h.viewers, userID)
	h.broadcastLocked()
}

func (h *Hub) disconnect(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	userID := h.clients[conn]
	delete(h.clients, conn)

	if userID != "" {
		delete(h.viewers, userID)
	}

	conn.Close()
	h.broadcastLocked()
}

func (h *Hub) broadcastLocked() {
	viewers := h.currentViewersLocked()

	for conn := range h.clients {
		err := conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			continue
		}

		err = conn.WriteJSON(viewers)
		if err != nil {
			userID := h.clients[conn]
			delete(h.clients, conn)

			if userID != "" {
				delete(h.viewers, userID)
			}

			conn.Close()
		}
	}
}

func (h *Hub) currentViewersLocked() []Viewer {
	viewers := make([]Viewer, 0, len(h.viewers))

	for _, viewer := range h.viewers {
		viewers = append(viewers, viewer)
	}

	sort.Slice(viewers, func(i int, j int) bool {
		return viewers[i].UserID < viewers[j].UserID
	})

	return viewers
}
