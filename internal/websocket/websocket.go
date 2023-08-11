package websocket

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error(err.Error())
		return ws, err
	}
	return ws, nil
}

func ServeWs(pool *Pool, w http.ResponseWriter, r *http.Request) {

	conn, err := Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &Client{
		Conn:       conn,
		Pool:       pool,
		RemoteAddr: conn.RemoteAddr().String(),
	}

	slog.Info("Client Connected", "client IP addr", client,
			"Previous length of clients", len(pool.Clients),
		)
	client.Read()
}
