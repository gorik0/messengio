package main

import (
	error2 "messengio/utils/error"
	"testing"
)

func TestStoreSql(t *testing.T) {
	store, err := NewSqlStore("./db/chat.db")
	error2.TestHandlerError(t, err, "error while creating new chat")
	testStore(t, store)
}

func testStore(t *testing.T, store Store) {

	//	:::ADD MSG

	msg := "HELO"
	err := store.PushMessage(msg)
	error2.TestHandlerError(t, err, "while pushing msg")

	//	::: RECEIVE MSG

	messages, err := store.GetMessages()
	error2.TestHandlerError(t, err, "while getting messages")
	if len(messages) == 0 {
		t.Error("no messages returned")
	}
}
