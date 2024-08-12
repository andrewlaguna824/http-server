package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"

	"http-server/server"
)

func main() {
	var port string
	flag.StringVar(&port, "port", "8080", "--port 8080")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConnectionString := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.New(ctx, dbConnectionString)
	if err != nil {
		log.Fatalf("Failed to create db connection pool %s", err.Error())
	}
	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("Failed to ping db %s", err.Error())
	}

	s := server.NewServer(port, pool)

	errChan := make(chan error)
	go func() {
		log.Printf("Running HTTP server on port %s", port)
		errChan <- s.Run(ctx)
	}()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err := <-errChan:
		log.Fatalf("API Server error: %s", err.Error())
	case <-signalChan:
		log.Println("Received signal interrupt")
		cancel()
	}
}
