package chat

import (
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	name          string
	conn          *websocket.Conn
	broadcastChan chan []byte
}

func (c *Client) SetName(name string) {
	c.name = name
}
func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn:          conn,
		broadcastChan: make(chan []byte),
	}
}

func (c *Client) listen() {
	go func() {
		for msg := range c.broadcastChan {
			c.conn.WriteMessage(websocket.TextMessage, msg)
		}
	}()
}
func (c *Client) write(msgBytes []byte) {

	log.Println("SENDING")
	log.Println(string(msgBytes))
	c.broadcastChan <- msgBytes
}
