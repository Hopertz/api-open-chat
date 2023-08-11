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

func (m MessageModel) FetchRoomMessages(roomID int, page, pageSize int) ([]Message, error) {

	offset := (page - 1) * pageSize

	query := `
        SELECT message.id, message.content , message.image_url , message.created_at , users.id
        FROM message
        INNER JOIN users ON users.id = message.user_id
        WHERE message.room_id = $1
        AND message.to_user_id = $2
        ORDER BY message.id DESC
        OFFSET $3
        LIMIT $4
    `

	rows, err := m.DB.Query(query, roomID, 0, offset, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		err := rows.Scan(
			&message.ID, &message.Content, &message.ImageUrl, &message.CreatedAt, &message.UserId,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// reverse the order of the messages
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}
