package websocket

import (
	"encoding/json"
	"github/hopertz/api-open-chat/internal/data"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn

	RemoteAddr string `json:"remote_addr"`

	Uid int `json:"uid"`

	RoomId int `json:"room_id"`

	Pool *Pool
}

const msgTypeJoin = 0

// const msgTypeLeave = 1
const msgTypeRoomChat = 2
const msgPrivateChat = 3

type MsgData struct {
	Uid      int       `json:"uid"`
	ToUid    int       `json:"to_uid"`
	Content  string    `json:"content"`
	ImageUrl string    `json:"image_url"`
	RoomId   int       `json:"room_id"`
	Time     time.Time `json:"time"`
}

type Message struct {
	Status int             `json:"status"`
	Data   MsgData         `json:"data"`
	Conn   *websocket.Conn `json:"conn"`
}

func (c *Client) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("Client Address", c.RemoteAddr),)
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {

		_, message, err := c.Conn.ReadMessage()

		if err != nil {
			slog.Error("Error reading message", "error", err)
			return
		}

		var msg Message

		if err := json.Unmarshal(message, &msg); err != nil {
			slog.Error("Error unmarshaling Json message", "error", err)
			continue
		}

		msgType := msg.Status

		msg.Conn = c.Conn
		c.RoomId = msg.Data.RoomId
		c.Uid = msg.Data.Uid

		if _, ok := c.Pool.Clients[c]; !ok {
			c.Pool.Clients[c] = msg.Data.Uid
		}

		switch {
		case msgType == msgTypeJoin:
			if msg.Data.RoomId > 0 {

				if room, ok := c.Pool.Rooms[msg.Data.RoomId]; ok {

					if !checkIfClientInroom(room, c) {
						c.Pool.Rooms[msg.Data.RoomId] = append(room, c)
						c.Pool.Broadcast <- msg
					}

				}
			}
		case msgType == msgTypeRoomChat:
			if msg.Data.RoomId != 0 {
				data := data.Message{
					UserId:    c.Uid,
					RoomId:    c.RoomId,
					Content:   msg.Data.Content,
					ImageUrl:  msg.Data.ImageUrl,
					CreatedAt: time.Now(),
				}

				err := c.Pool.model.MessageModel.Insert(data)

				if err != nil {
					slog.Error(err.Error())
				}

				c.Pool.Broadcast <- msg
			}

		case msgType == msgPrivateChat:
			if msg.Data.ToUid != 0 {
				data := data.Message{
					UserId:    c.Uid,
					ToUserId:  msg.Data.ToUid,
					Content:   msg.Data.Content,
					ImageUrl:  msg.Data.ImageUrl,
					CreatedAt: time.Now(),
				}

				err := c.Pool.model.MessageModel.Insert(data)

				if err != nil {
					slog.Error("Error inserting message to the database", "error", err)
				}

				for cli, uid := range c.Pool.Clients {
					if uid == msg.Data.ToUid {
						if err := cli.Conn.WriteJSON(msg); err != nil {
							slog.Error("Error sending private msg to client", "error", err)
							return
						}
					}
				}
			}
		}

	}
}
