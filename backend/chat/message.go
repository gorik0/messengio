package chat

import (
	"errors"
	"fmt"
	"strconv"
)

const MESSAGE_TYPE_TEXT = 1
const MESSAGE_TYPE_NAME = 2
const MESSAGE_TYPE_LEAVE = 3
const MESSAGE_TYPE_TYPING = 4
const MESSAGE_TYPE_STOP_TYPING = 5

type Message struct {
	msgType int
	author  *Client
	room    *Room
	content []byte
}

func Encode(m *Message) []byte {

	var content string
	if m.msgType == MESSAGE_TYPE_TEXT {
		content = string(m.content)

	}

	output := fmt.Sprintf("%d%s | %s", m.msgType, m.author.name, content)
	return []byte(output)

}

func parseMsgData(data []byte) (int, []byte, error) {
	//	::: CHECK msgTYPE
	println(string(data[1:]))
	msgType, err := strconv.Atoi(string(data[0]))
	if err != nil {
		return 0, nil, errors.New("PARSE!!!")
	}

	if msgType > 5 || msgType < 1 {
		return 0, nil, errors.New("Error msg type!!!")
	}
	return msgType, data[1:], nil
}

func NewMessage(typ int, client *Client, room *Room, content []byte) *Message {
	return &Message{
		msgType: typ,
		author:  client,
		room:    room,
		content: content,
	}
}
