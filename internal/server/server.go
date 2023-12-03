package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	router *http.ServeMux
}

func NewServer() *Server {
	return &Server{
		router: http.NewServeMux(),
	}
}

func (s *Server) Router() *http.ServeMux {
	s.router.HandleFunc("/ws", s.handleWebSocket)
	return s.router
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	go s.handleConnection(conn)
}

func (s *Server) handleConnection(conn *websocket.Conn) {
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		log.Printf("Received message: %s\n", message)

		err = conn.WriteMessage(messageType, []byte(fmt.Sprintf("Server received message: %s", message)))
		if err != nil {
			log.Println(err)
			break
		}
	}
}

func (s *Server) Close() {
	// Perform any necessary cleanup here
}
