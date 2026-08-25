package main

import (
	"context"
	"log"
	"time"

	"github.com/arafat2020/go-phinix/internal/executor"
)

func main() {
	cfg := executor.Config{
		DSN:          "postgres://postgres:postgrespassword@localhost:5432/phinix?sslmode=disable",
		MaxOpenConns: 10,
	}

	ex, err := executor.New(cfg)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer ex.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := ex.Ping(ctx); err != nil {
		log.Fatalf("ping: %v", err)
	}

	log.Println("connected to postgres ✔")
}
