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

		c.Pool.Clients[c] = true // add client to pool

		fmt.Println(c.Pool.Clients)
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
				if room, ok := c.Pool.Rooms[msg.Data.RoomId]; ok {

					if !checkIfClientInroom(room, c) {
						c.Pool.Rooms[msg.Data.RoomId] = append(room, c)
						c.Pool.Broadcast <- msg
					}

				}
			}
		} else if msgType == msgTypeRoom {
			if msg.Data.RoomId != "" {
				if room, ok := c.Pool.Rooms[msg.Data.RoomId]; ok {
					for _, receiver := range room {
						if err := receiver.Conn.WriteJSON(msg); err != nil {
							log.Println(err)
							return
						}
					}
				}
			}
		}

	}
}
