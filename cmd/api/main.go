package main

import (
	"context"
	"database/sql"
	"flag"
	"github/hopertz/api-open-chat/internal/data"
	"github/hopertz/api-open-chat/internal/websocket"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"

	"log"
	"log/slog"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

type application struct {
	config config
	pool   *websocket.Pool
	wg     sync.WaitGroup
	models data.Models
}

func init() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)

}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4141, "API server port")
	flag.StringVar(&cfg.env, "env", "production", "Environment (development|Staging|production")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("CHAT_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max ilde connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection  connections")

	flag.Parse()

	conn, err := openDB(cfg)

	if err != nil {
		log.Fatal(err, nil)
	}

	defer conn.Close()

	dbConn := data.NewModels(conn)
	pool := websocket.NewPool(dbConn)

	app := &application{
		config: cfg,
		pool:   pool,
		models: dbConn,
	}

	go pool.Start()

	slog.Info("Starting server on port", "port", cfg.port)

	err = app.serve()

	if err != nil {
		log.Fatal(err, nil)
	}

}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	return db, nil
}
