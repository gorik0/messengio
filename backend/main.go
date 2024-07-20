package main

import (
	"github.com/gorilla/websocket"
	"log"
	. "messengio/utils/error"
	"net/http"
)

func main() {

	server := NewServer()
	go server.listenChannels()
	port := ":8080"
	log.Println("Starting server at ::: ", port)
	HandlerError(http.ListenAndServe(port, server), "error while starting server")
}

type Server struct {
	upgrader         websocket.Upgrader
	mux              *http.ServeMux
	client           []*Client
	msgChannel       chan Message
	registerClient   chan *Client
	unregisterClient chan *Client
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	conn, err := s.upgrader.Upgrade(w, r, nil)
	HandlerError(err, "while upgrading to websocket")

	//defer conn.Close()

	newClient := NewClient(conn)
	s.registerClient <- newClient

	go readMsgs(newClient, s)
}

func (s *Server) listenChannels() {
	for {

		select {
		case client := <-s.unregisterClient:
			{
				for i, c := range s.client {
					if c == client {
						s.client = append(s.client[:i], s.client[i:]...)
						break
					}

				}
			}
		case client := <-s.registerClient:
			{

				s.client = append(s.client, client)
			}
		case msg := <-s.msgChannel:
			{
				for _, client := range s.client {

					if msg.author.name == "" {
						msg.author.name = "bobi"
					}
					msgUI := msg.author.name + " | " + string(msg.content)
					client.write([]byte(msgUI))

				}
			}
		}
	}
}

func readMsgs(client *Client, server *Server) {
	for {

		//	::RECIVE MSG
		_, msgBytes, err := client.conn.ReadMessage()
		if err != nil {
			HandlerErrorLite(err, "error while reading message")
			return
		}

		log.Println("GOTTA msg ", string(msgBytes)) //	:::SEND MSG TO ALL

		msg := Message{
			content: msgBytes,
			author:  client,
		}
		server.msgChannel <- msg

	}
}

func NewServer() *Server {

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	msgChannel := make(chan Message)
	registerClient := make(chan *Client)
	unregisterClient := make(chan *Client)

	return &Server{
		upgrader:         upgrader,
		msgChannel:       msgChannel,
		registerClient:   registerClient,
		unregisterClient: unregisterClient,
	}
}
