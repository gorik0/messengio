package main

type Store interface {
	GetMessages() ([]string, error)
	PushMessage(msg string) error
}

type MemoryStore struct {
	messages []string
}

func (m *MemoryStore) GetMessages() ([]string, error) {
	return m.messages, nil
}

func (m *MemoryStore) PushMessage(msg string) error {

	m.messages = append(m.messages, msg)
	return nil
}

var _ Store = &MemoryStore{}

func NewStore() Store {
	return &MemoryStore{}
}
