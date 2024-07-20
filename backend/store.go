package main

import (
	"database/sql"
)

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

type SqlStore struct {
	db *sql.DB
}

func (s SqlStore) GetMessages() ([]string, error) {
	q := "SELECT message FROM messages  "
	query, err := s.db.Query(q, nil)
	var msgs []string
	if err != nil {
		return msgs, err
	}

	for query.Next() {
		var msg string
		err := query.Scan(&msg)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil

}

func (s SqlStore) PushMessage(msg string) (err error) {
	q := "INSERT INTO messages (message) VALUES (?)"
	_, err = s.db.Exec(q, msg)
	return err
}

var _ Store = &SqlStore{}

func NewSqlStore(db *sql.DB) (Store, error) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS messages (message TEXT)")
	if err != nil {
		return nil, err
	}

	return &SqlStore{db: db}, nil

}
