package main

import (
	"context"
	"fmt"
	"log"

	"github.com/primalcs/words_of_wisdom/internal/server"
)

func main() {
	log.Println("starting server...")

	ctx := context.Background()
	cfg, err := server.LoadConfig()
	if err != nil {
		panic(err)
	}

	server, err := server.NewServer(ctx, cfg)
	if err != nil {
		panic(fmt.Errorf("error when creating a server: %v", err))
	}

	defer func() {
		if err := server.Close(ctx); err != nil {
			panic(err)
		}
	}()

	if err = server.Run(ctx); err != nil {
		log.Printf("server error: %v", err)
	}
}
