package websocket

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn

	RemoteAddr string `json:"remote_addr"`

	Uid string `json:"uid"`

	Username string `json:"username"`

	RoomId string `json:"room_id"`

	Pool *Pool
}

const msgTypeJoin = 0
const msgTypeLeave = 1
const msgTypeRoom = 2
const msgPrivate = 3

type MsgData struct {
	Uid      string        `json:"uid"`
	Username string        `json:"username"`
	AvatarId string        `json:"avatar_id"`
	ToUid    string        `json:"to_uid"`
	Content  string        `json:"content"`
	ImageUrl string        `json:"image_url"`
	RoomId   string        `json:"room_id"`
	Count    int           `json:"count"`
	List     []interface{} `json:"list"`
	Time     int64         `json:"time"`
}

type Message struct {
	Status int             `json:"status"`
	Data   MsgData         `json:"data"`
	Conn   *websocket.Conn `json:"conn"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {

		fmt.Println(c.Pool.Rooms)
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		msgType := msg.Status

		c.RoomId = msg.Data.RoomId
		c.Username = msg.Data.Username
		c.Uid = msg.Data.Uid

		if msgType == msgTypeJoin {

			if msg.Data.RoomId != "" {
				if _, ok := c.Pool.Rooms[msg.Data.RoomId]; ok {
					c.Pool.Rooms[msg.Data.RoomId] = append(c.Pool.Rooms[msg.Data.RoomId], c)
					c.Pool.Clients[c] = true
				}

				c.Pool.Broadcast <- msg

			}
		}

	}
}
