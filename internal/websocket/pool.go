package websocket

import (
	"fmt"
	"github/hopertz/api-open-chat/internal/data"

	"log/slog"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]int
	Broadcast  chan Message
	Rooms      map[int][]*Client
	model      data.Models
}

func NewPool(model data.Models) *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]int),
		Broadcast:  make(chan Message),
		Rooms: map[int][]*Client{
			1: {},
			2: {},
			3: {},
			4: {},
			5: {},
			6: {},
		},

		model: model,
	}
}
func (pool *Pool) Start() {
	for {
		select {

		case client := <-pool.Unregister:
			delete(pool.Clients, client)

			for roomId, receivers := range pool.Rooms {
				for i, receiver := range receivers {
					if receiver == client {
						// remove the element at index i from a.
						pool.Rooms[roomId] = append(pool.Rooms[roomId][:i], pool.Rooms[roomId][i+1:]...)
						break
					}

				}
			}

			slog.Info("User left the connection: ", "id", client.Uid)
			slog.Info("Pool size Remaining", "length", len(pool.Clients))

		case message := <-pool.Broadcast:
			var msg map[int]string
			if message.Status == 0 {
				msg = map[int]string{
					message.Data.Uid: fmt.Sprintf("Has joined in room %d: ", message.Data.RoomId),
				}
			} else {
				msg = map[int]string{
					message.Data.Uid: message.Data.Content,
				}
			}
			if room, ok := pool.Rooms[message.Data.RoomId]; ok {
				for _, receiver := range room {
					if message.Conn != receiver.Conn {
						if err := receiver.Conn.WriteJSON(msg); err != nil {
							slog.Error("Error writing message to client in room", slog.Group("room", message.Data.RoomId, "client", receiver.Uid, "error", err))
							continue
						}
					}
				}
			}

		}
	}
}
