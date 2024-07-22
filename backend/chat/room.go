package chat

type Room struct {
	ID      string
	clients []*Client
}

type Instruction struct {
	room   *Room
	client *Client
}

func NewInstruction(room *Room, client *Client) *Instruction {
	return &Instruction{
		room:   room,
		client: client,
	}
}

func (r *Room) ClientCount() int {
	return len(r.clients)

}
