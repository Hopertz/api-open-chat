package websocket

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

// ints on Body Type
// 0  for new user
// 1  for normal text
type Body struct {
	Type int    `json:"type"`
	Text string `json:"text"`
	User string `json:"user"`
}

type Message struct {
	Type  int      `json:"type"`
	Body  Body     `json:"body"`
	Users []string `json:"users"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		var data Body

		if err := json.Unmarshal(p, &data); err != nil {
			continue
		}

		message := Message{Type: messageType, Body: data}

		if message.Body.Type == 0 {
			c.Pool.Users = append(c.Pool.Users, message.Body.User)
		}

		message.Users = c.Pool.Users

		c.Pool.Broadcast <- message

		fmt.Printf("Message Received: %+v\n", message)
	}
}
