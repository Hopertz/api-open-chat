package websocket

import log "github.com/sirupsen/logrus"

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
	Users      []string
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
		Users:      []string{},
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			log.Info("Size of Connection Pool: ", len(pool.Clients))
			log.Info("New user joined")

		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			log.Info("Size of Connection Pool: ", len(pool.Clients))
			log.Info(" Another One Bites the Dust Song by Queen")

		case message := <-pool.Broadcast:
			log.Info("Sending message to all clients in Pool")
			for client := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					log.Info(err)
					return
				}
			}
		}
	}
}
