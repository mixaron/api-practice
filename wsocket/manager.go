package wsocket

import (
	"github.com/gofiber/websocket/v2"
	"sync"
)

type Manager struct {
	conns map[*websocket.Conn]bool
	mu    sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (m *Manager) Add(conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.conns[conn] = true
}

func (m *Manager) Remove(conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.conns, conn)
}

func (m *Manager) List() []*websocket.Conn {
	m.mu.RLock()
	defer m.mu.RUnlock()

	conns := make([]*websocket.Conn, 0, len(m.conns))
	for conn := range m.conns {
		conns = append(conns, conn)
	}
	return conns
}
