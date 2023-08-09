package websocket

import (
	"fmt"
	"github/hopertz/api-open-chat/internal/data"

	log "github.com/sirupsen/logrus"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
	Rooms      map[int][]*Client
	model      data.Models
}

func NewPool(model data.Models) *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
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

			log.Info("Size of Connection Pool: ", len(pool.Clients))
			fmt.Println(pool.Clients)
			log.Info(" Another One Bites the Dust")

		case message := <-pool.Broadcast:

			if message.Data.RoomId != 0 {
				log.Info(message.Data.Uid, " Has joined in room: ", message.Data.RoomId)
				if receivers, ok := pool.Rooms[message.Data.RoomId]; ok {
					for _, receiver := range receivers {
						if err := receiver.Conn.WriteJSON(message); err != nil {
							log.Info(err)
							return
						}
					}
				}
			}
		}
	}
}
