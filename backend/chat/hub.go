package chat

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const EMPTY_ROOM_TIMEOUT = 1 * time.Minute

type Hub interface {
	Register(room *Room, client *Client)
	Unregister(room *Room, client *Client)
	Broadcast(msg *Message)

	handleRegister(room *Room, client *Client)
	handleUnregister(room *Room, client *Client)
	handleBroadcast(msg *Message)

	GetRoom(id string) *Room
	CreateRoom() *Room
	Run()
	ListenClient(client *Client, room *Room)
}

type HubIts struct {
	broadcastChan chan *Message
	register      chan *Instruction
	unregister    chan *Instruction
	rooms         map[string]*Room
}

func (h *HubIts) Register(room *Room, client *Client) {
	h.register <- NewInstruction(room, client)
}

func (h *HubIts) Unregister(room *Room, client *Client) {
	h.unregister <- NewInstruction(room, client)
}

func (h *HubIts) Broadcast(msg *Message) {
	h.broadcastChan <- msg
}

func (h *HubIts) handleRegister(room *Room, client *Client) {
	room.clients = append(room.clients, client)
	h.rooms[room.ID] = room
	client.listen()

}

func (h *HubIts) handleUnregister(room *Room, client *Client) {

	//	::: Delete client from room
	for i, clientInRoom := range room.clients {
		if clientInRoom == client {
			room.clients = append(room.clients[:i], room.clients[i+1:]...)
			//	:::: Create LEAVE_type message

		}
	}
	if client.name != "" {
		content := fmt.Sprintf("Client %s has left from room", client.name)
		msg := NewMessage(MESSAGE_TYPE_LEAVE, client, room, []byte(content))
		go h.Broadcast(msg)
	}
	h.scheduleRoomTermination(room)

}

func (h *HubIts) handleBroadcast(msg *Message) {
	bytes := Encode(msg)
	for _, client := range msg.room.clients {
		log.Println(msg.author.name)

		if client.name != msg.author.name {
			client.write(bytes)
		}
	}
}

func (h *HubIts) GetRoom(id string) *Room {
	if room, ok := h.rooms[id]; ok {
		return room
	} else {
		return nil
	}
}

func (h *HubIts) CreateRoom() *Room {
	room := &Room{
		ID:      uuid.New().String(),
		clients: make([]*Client, 0),
	}
	h.rooms[room.ID] = room
	return room
}

func (h *HubIts) Run() {
	go func() {

		for {
			select {
			case register := <-h.register:

				h.handleRegister(register.room, register.client)

			case unregister := <-h.unregister:

				h.handleUnregister(unregister.room, unregister.client)

			case broadcast := <-h.broadcastChan:

				h.handleBroadcast(broadcast)

			}

		}
	}()

}

func (h *HubIts) ListenClient(client *Client, room *Room) {
	for {
		messageType, mesg, err := client.conn.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			log.Println("Invalid msg type reciebver . Disconencting...", err.Error())
			h.Unregister(room, client)
			return
		}
		mesgType, bytes, err := parseMsgData(mesg)
		if err != nil {
			log.Println("Invalid to parse MESG . Disconencting...")
			h.Unregister(room, client)
			return

		}
		msg := NewMessage(mesgType, client, room, bytes)
		if mesgType == MESSAGE_TYPE_NAME {
			client.SetName(string(msg.content))
			continue
		}
		if client.name == "" && msg.msgType != MESSAGE_TYPE_NAME {
			log.Println("Unknown client . Disconencting...")
			h.Unregister(room, client)
			return

		}
		h.Broadcast(msg)

	}
}

var _ Hub = &HubIts{}

func NewHubIts() Hub {
	return &HubIts{
		broadcastChan: make(chan *Message),
		register:      make(chan *Instruction),
		unregister:    make(chan *Instruction),
		rooms:         make(map[string]*Room),
	}
}

func (h *HubIts) scheduleRoomTermination(room *Room) {
	go func() {
		time.AfterFunc(EMPTY_ROOM_TIMEOUT, func() {
			if room.ClientCount() > 0 || h.rooms[room.ID] == nil {
				return
			}
			delete(h.rooms, room.ID)
		})
	}()
}
