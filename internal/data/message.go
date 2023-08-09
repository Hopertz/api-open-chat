package data

import (
	"database/sql"
	"time"
)

type Message struct {
	ID        int
	UserId    int
	ToUserId  int
	RoomId    int
	Content   string
	ImageUrl  string
	CreatedAt time.Time
}

type MessageModel struct {
	DB *sql.DB
}

func (m MessageModel) Insert(msg Message) error {
	stmt := `INSERT INTO message (user_id, to_user_id, room_id, content, image_url, created_at)
	VALUES ($1, $2, $3, $4, $5, $6)`

	args := []interface{}{msg.UserId, msg.ToUserId, msg.RoomId, msg.Content, msg.ImageUrl, msg.CreatedAt}
	_, err := m.DB.Exec(stmt, args...)
	if err != nil {
		return err
	}

	return nil
}
