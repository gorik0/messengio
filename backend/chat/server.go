package chat

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
)

type Server struct {
	upgrader websocket.Upgrader
	mux      *http.ServeMux
	hub      Hub
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) Run(port string) {
	log.Println("Start server on port ::: ", port)
	log.Fatal(http.ListenAndServe(":"+port, s))

}
func (s *Server) configureRoutes() {
	s.mux.HandleFunc("OPTIONS /*", s.handleOptions)
	s.mux.HandleFunc("GET /room/{room}", s.handleRoomGet)
	s.mux.HandleFunc("POST /room", s.handleRoomCreate)
	s.mux.HandleFunc("GET /ws/{room}", s.handleWS)

}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {

	//	:::UPGRADE UPGRADR
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	//	:::CREATE  CLIENT
	client := NewClient(conn)
	//	:::REGISTER  CLIENT
	room := s.GetRoomFromUrlPath(w, r)
	if room == nil {
		return
	}
	s.hub.Register(room, client)

	//	:::LISTEN

	s.hub.ListenClient(client, room)
}

func (s *Server) handleRoomCreate(w http.ResponseWriter, r *http.Request) {
	println("creating ... rooom ....")
	s.setCORSpolicy(w)
	room := s.hub.CreateRoom()
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(room.ID))
}

func (s *Server) handleRoomGet(w http.ResponseWriter, r *http.Request) {
	s.setCORSpolicy(w)

	s.GetRoomFromUrlPath(w, r)
}

func (s *Server) GetRoomFromUrlPath(w http.ResponseWriter, r *http.Request) *Room {
	roomId := r.PathValue("room")
	if roomId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}

	if room := s.hub.GetRoom(roomId); room == nil {
		w.WriteHeader(http.StatusBadRequest)
		return room
	} else {
		return room
	}

}

func (s *Server) handleOptions(w http.ResponseWriter, r *http.Request) {
	s.setCORSpolicy(w)
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) setCORSpolicy(writer http.ResponseWriter) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
}

func NewServer(hub Hub) *Server {
	//::: setup CORS

	var corsAllowOrigin string

	if os.Getenv("CORS_ALLOW_ORIGIN") != "" {
		corsAllowOrigin = os.Getenv("CORS_ALLOW_ORIGIN")
	} else {
		corsAllowOrigin = "*"
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  0,
		WriteBufferSize: 0,
		CheckOrigin: func(r *http.Request) bool {
			if corsAllowOrigin == "*" {
				return true
			} else {
				origin := r.Header.Get("Origin")

				return origin == corsAllowOrigin
			}
		},
	}
	mux := http.NewServeMux()
	server := Server{
		upgrader: upgrader,
		mux:      mux,
		hub:      hub,
	}

	server.configureRoutes()

	return &server

}
