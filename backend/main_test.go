package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"

	"messengio/utils/constant"
	. "messengio/utils/error"
	"net/http/httptest"
	"testing"
)

func TestWEbsocketConnection(t *testing.T) {

	store := NewStore()
	server := NewServer(store)
	s := httptest.NewServer(server)
	defer s.Close()

	dialer := websocket.Dialer{}
	dial, response, err := dialer.Dial("ws"+s.URL[4:]+"/ws", nil)
	defer dial.Close()
	TestHandlerError(t, err, "create dial ws")

	if response.StatusCode != http.StatusSwitchingProtocols {
		t.Errorf("wrong http status, expected %d, got %d", http.StatusSwitchingProtocols, response.StatusCode)
	}

}
func TestWEbsocketMessage(t *testing.T) {

	store, err := NewSqlStore(constant.DATABASE_PATH)
	TestHandlerError(t, err, "while creating store sql")
	server := NewServer(store)
	s := httptest.NewServer(server)
	defer s.Close()

	dialer := websocket.Dialer{}
	dial, _, err := dialer.Dial("ws"+s.URL[4:]+"/ws", nil)
	defer dial.Close()
	TestHandlerError(t, err, "create dial ws")

	//:::: WRITE MSG
	TestHandlerError(t, dial.WriteMessage(websocket.TextMessage, []byte("hello world")), "error while writing msg")
	//	:::: GET msg from database
	time.Sleep(time.Second * 1)

	messages, err := store.GetMessages()
	if err != nil {
		TestHandlerError(t, err, "create dial ws")
	}
	if len(messages) != 1 {
		t.Errorf("wrong messages, expected %d, got %d", 1, len(messages))
	}
}
