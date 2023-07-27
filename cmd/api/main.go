package main

import (
	"flag"
	"github/hopertz/api-open-chat/internal/websocket"
	"sync"

	log "github.com/sirupsen/logrus"
)

type config struct {
	port int
	env  string
}

type application struct {
	config config
	pool   *websocket.Pool
	wg     sync.WaitGroup
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4141, "API server port")
	flag.StringVar(&cfg.env, "env", "production", "Environment (development|Staging|production")

	flag.Parse()

	pool := websocket.NewPool()

	app := &application{
		config: cfg,
		pool:   pool,
	}

	go pool.Start()

	err := app.serve()
	if err != nil {
		log.Fatal(err, nil)
	}

}
