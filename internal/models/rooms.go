package models

import (
	"database/sql"
	"time"
)

type RoomModelInterface interface {
	Latest(id int) ([]*Message, error)
}

type Room struct {
	ID       int
	Name     string
	Members  []*Member
	Messages []*Message
}

type Member struct {
	UserID int
	RoomID int
}

type Message struct {
	ID        int
	UserID    int
	RoomID    int
	Text      string
	Timestamp time.Time
}

type RoomModel struct {
	DB *sql.DB
}

func (m *RoomModel) Latest(id int) ([]*Message, error) {
	stmt := `SELECT message_id, user_id, room_id, message_txt, timestamp
	FROM messages
	Where id = ?
	ORDER By timestamp DESC
	LIMIT 100`

	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	messages := []*Message{}

	for rows.Next() {
		msg := &Message{}

		err = rows.Scan(&msg.ID, &msg.UserID, &msg.RoomID, &msg.Text, &msg.Timestamp)
		if err != nil {
			return nil, err
		}

		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
