package main

import "github.com/gorilla/websocket"

type Client struct {
	conn *websocket.Conn
	name string
}

func (c *Client) write(mes []byte) (err error) {

	return c.conn.WriteMessage(websocket.TextMessage, mes)
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{conn: conn}
}
