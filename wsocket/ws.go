package wsocket

import (
	"api-practice/internal/service"
	"encoding/json"
	"fmt"
	"github.com/gofiber/websocket/v2"
)

type Server struct {
	conns          *Manager
	userService    service.UserService
	articleService service.ArticleService
}

func NewServer(u service.UserService, a service.ArticleService, m *Manager) *Server {
	return &Server{
		conns:          m,
		userService:    u,
		articleService: a,
	}
}

func (s *Server) SendMessage(value string) {
	s.Broadcast([]byte(value))
}

func (s *Server) HandleWS(c *websocket.Conn, userID string) {
	fmt.Println("new connection from user: ", userID)
	s.conns.Add(c)
	s.readLoop(c, userID)
}

func (s *Server) readLoop(conn *websocket.Conn, userID string) {
	defer func() {
		conn.Close()
		s.conns.Remove(conn)
		_ = s.userService.SetUserLastOnlineTime(userID)
	}()

	for {
		mt, _, err := conn.ReadMessage()

		if err := conn.WriteMessage(mt, []byte("connect successful")); err != nil {
			fmt.Println("failed to connect user:", err)
			break
		}

		userLastOnline, err := s.userService.GetUserLastOnlineTime(userID)
		if err != nil {
			fmt.Println(err)
			break
		}

		articles, _ := s.articleService.GetAllArticlesAfterTime(userLastOnline)

		if len(articles) > 0 {
			if err := conn.WriteMessage(mt, []byte("new messages:")); err != nil {
				fmt.Println("failed to send message notification:", err)
				return
			}

			for _, article := range articles {
				jsonArticle, err := json.Marshal(article)
				if err != nil {
					fmt.Println("failed to marshal article:", err)
					continue
				}
				if err := conn.WriteMessage(mt, jsonArticle); err != nil {
					fmt.Println("failed to send article:", err)
					break
				}
			}
		}

		fmt.Println("user connected: " + userID)
	}
}

func (s *Server) Broadcast(b []byte) {
	for _, conn := range s.conns.List() {
		go func(conn *websocket.Conn) {
			if err := conn.WriteMessage(websocket.TextMessage, b); err != nil {
				fmt.Println("write error: ", err)
			}
		}(conn)
	}
}
