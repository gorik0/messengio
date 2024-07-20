package main

import (
	"github.com/gorilla/websocket"
	. "messengio/utils/error"
	"net/http"
)

type Server struct {
	upgrader websocket.Upgrader
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.handleWs)
	mux.ServeHTTP(w, r)

}

func (s *Server) handleWs(writer http.ResponseWriter, request *http.Request) {

	conn, err := s.upgrader.Upgrade(writer, request, nil)
	HandlerError(err, "while upgrading to websocket")
	defer conn.Close()
}

func main() {

	server := NewServer()
	http.ListenAndServe(":8080", server)
}

func NewServer() *Server {

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	return &Server{upgrader: upgrader}
}
