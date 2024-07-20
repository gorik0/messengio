package main

import (
	"github.com/gorilla/websocket"
	. "messengio/utils/error"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWEbsocketConnection(t *testing.T) {

	server := NewServer()
	s := httptest.NewServer(server)
	defer s.Close()

	println(s.URL[4:])

	dialer := websocket.Dialer{}
	dial, response, err := dialer.Dial("ws"+s.URL[4:]+"/ws", nil)
	defer dial.Close()
	TestHandlerError(t, err, "create dial ws")

	if response.StatusCode != http.StatusSwitchingProtocols {
		t.Errorf("wrong http status, expected %d, got %d", http.StatusSwitchingProtocols, response.StatusCode)
	}

}
func TestWEbsocketMessage(t *testing.T) {

	server := NewServer()
	s := httptest.NewServer(server)
	defer s.Close()

	println(s.URL[4:])

	dialer := websocket.Dialer{}
	dial, _, err := dialer.Dial("ws"+s.URL[4:]+"/ws", nil)
	defer dial.Close()
	TestHandlerError(t, err, "create dial ws")

	TestHandlerError(t, dial.WriteMessage(websocket.TextMessage, []byte("hello world")), "error while writing msg")

}
