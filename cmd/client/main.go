package main

import (
	"context"
	"log"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/primalcs/words_of_wisdom/internal/client"
)

func main() {
	log.Println("starting client...")
	ctx := context.Background()

	cfg, err := client.LoadConfig()
	if err != nil {
		panic(err)
	}

	clientInstance, err := client.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := clientInstance.Close(); err != nil {
			panic(err)
		}
	}()

	for {
		log.Println("running client...")
		err = clientInstance.HandleConnection(ctx)
		if err != nil {
			log.Printf("client error: %v", err)
		}

		time.Sleep(3 * time.Second)
	}
}
