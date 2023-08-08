package data

import (
	"database/sql"
	"time"
)

type Message struct {
	ID        string
	UserId    string
	ToUserId  string
	RoomId    int
	Content   string
	ImageUrl  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MessageModel struct {
	DB *sql.DB
}

func (m MessageModel) Insert(userId int, toUserId int, roomId int, content string, imageUrl string) (string, error) {
	stmt := `INSERT INTO messages (user_id, to_user_id, room_id, content, image_url, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	var id string
	err := m.DB.QueryRow(stmt, userId, toUserId, roomId, content, imageUrl, time.Now(), time.Now()).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}
