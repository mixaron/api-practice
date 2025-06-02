package wsocket

import (
	"fmt"
	"github.com/gofiber/websocket/v2"
	"time"
)

type Server struct {
	conns      map[*websocket.Conn]bool
	userMap    map[string]time.Time
	messageMap map[string]time.Time
}

type MessageRequest struct {
	Message string `json:"message"`
}

func NewServer() *Server {
	return &Server{
		conns:      make(map[*websocket.Conn]bool),
		userMap:    make(map[string]time.Time),
		messageMap: make(map[string]time.Time),
	}
}

func (s *Server) SendMessage(value string) {
	s.messageMap[value] = time.Now()
	s.Broadcast([]byte(value))
}

func (s *Server) HandleWS(c *websocket.Conn, userID string) {
	fmt.Println("new connection from user: ", userID)

	s.conns[c] = true

	s.readLoop(c, userID)
}

func (s *Server) readLoop(conn *websocket.Conn, userID string) {
	defer func() {
		err := conn.Close()
		if err != nil {
			return
		}
		delete(s.conns, conn)
		s.userMap[userID] = time.Now()
	}()

	for {
		mt, _, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read error:", err)
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				fmt.Println("current user disc: " + userID)
				s.userMap[userID] = time.Now()
			}
			break
		}

		if err := conn.WriteMessage(mt, []byte("connect successful")); err != nil {
			fmt.Println("failed to connect user:", err)
			break
		}

		if _, exists := s.userMap[userID]; !exists {
			if err := conn.WriteMessage(mt, []byte("new messages:")); err != nil {
				fmt.Println("failed to send message history:", err)
				break
			}

			for message := range s.messageMap {
				if err := conn.WriteMessage(mt, []byte(message)); err != nil {
					fmt.Println("failed to send message history:", err)
					break
				}
			}
		} else {
			for key, value := range s.messageMap {
				if value.After(s.userMap[userID]) {
					if err := conn.WriteMessage(mt, []byte(key)); err != nil {
						fmt.Printf("failed to send message to %s: %v\n", userID, err)
						break
					}
				}
			}
		}

		fmt.Println("user connected: " + userID)
	}
}

func (s *Server) Broadcast(b []byte) {
	for conn := range s.conns {
		go func(conn *websocket.Conn) {
			if err := conn.WriteMessage(websocket.TextMessage, b); err != nil {
				fmt.Println("write error: ", err)
			}
		}(conn)
	}
}
