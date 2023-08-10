package websocket

import (
	"encoding/json"
	"fmt"
	"github/hopertz/api-open-chat/internal/data"
	"log"
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
const msgTypeLeave = 1
const msgTypeRoom = 2
const msgPrivate = 3

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
			log.Println(err)
			continue
		}

		msgType := msg.Status

		c.RoomId = msg.Data.RoomId
		c.Uid = msg.Data.Uid

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
		case msgType == msgTypeRoom:
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
					log.Println(err)
				}

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
