package main

import (
	"github.com/gorilla/websocket"
	. "messengio/utils/error"
	"net/http"
)

func main() {

	store, err := NewSqlStore("./db/chat.db")
	HandlerError(err, "while connecting to database")
	server := NewServer(store)
	http.ListenAndServe(":8080", server)
}

type Server struct {
	upgrader websocket.Upgrader
	mux      *http.ServeMux
	store    Store
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	s.mux.ServeHTTP(w, r)

}

func (s *Server) handleWs(writer http.ResponseWriter, request *http.Request) {

	conn, err := s.upgrader.Upgrade(writer, request, nil)
	HandlerError(err, "while upgrading to websocket")
	defer conn.Close()

	readMsgs(conn, s.store)
}

func readMsgs(conn *websocket.Conn, store Store) {
	for {
		//	 :::GET MSG
		_, msgBytes, err := conn.ReadMessage()

		if err != nil {
			HandlerErrorLite(err, "add message to store")
			break

		}
		println("store", store)
		println("GORIOoooo   :::: ", (string(msgBytes)))

		//	 :::ADD MSG
		err = store.PushMessage(string(msgBytes))
		if err != nil {
			HandlerErrorLite(err, "add message to store")
			break

		}

	}
}

func NewServer(store Store) *Server {
	mux := http.NewServeMux()

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	server := &Server{
		upgrader: upgrader,
		mux:      mux,
		store:    store,
	}
	mux.HandleFunc("/ws", server.handleWs)
	return server
}
