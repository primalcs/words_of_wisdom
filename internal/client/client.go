package client

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"

	"github.com/primalcs/words_of_wisdom/internal/hashcash"
	"github.com/primalcs/words_of_wisdom/pkg/protocol"
)

type Client struct {
	conn net.Conn

	hashCashMaxIterationsAmount int64
}

func NewClient(cfg *Config) (*Client, error) {
	conn, err := net.Dial("tcp", cfg.ServerAddress)
	if err != nil {
		return nil, err
	}
	log.Printf("connected to %s", cfg.ServerAddress)

	return &Client{
		conn:                        conn,
		hashCashMaxIterationsAmount: cfg.MaxIterations,
	}, nil
}

func (client *Client) Close() error {
	return client.conn.Close()
}

func (client *Client) HandleConnection(ctx context.Context) error {
	connBufReader := bufio.NewReader(client.conn)
	connBufWriter := bufio.NewWriter(client.conn)

	// request challenge
	if err := protocol.RequestChallenge(connBufWriter); err != nil {
		return fmt.Errorf("sending request for challenge failed, error: %w", err)
	}

	challenge := &hashcash.Challenge{}
	// receive challenge
	powPuzzle, err := protocol.ReceiveChallenge(connBufReader, challenge)
	if err != nil {
		return fmt.Errorf("reading challenge msg failed, error: %w", err)
	}
	log.Printf("solving puzzle: %v", powPuzzle)

	// solve challenge
	err = powPuzzle.Solve(client.hashCashMaxIterationsAmount)
	if err != nil {
		return fmt.Errorf("puzzle solving failed, error: %w", err)
	}
	log.Printf("puzzle solved: %v", powPuzzle)

	// send challenge solution
	err = protocol.SendChallengeSolution(connBufWriter, powPuzzle)
	if err != nil {
		return fmt.Errorf("sending request failed, error: %w", err)
	}
	log.Println("challenge solution sent to server")

	// receive resource
	resource, err := protocol.ReceiveResource(connBufReader)
	if err != nil {
		return fmt.Errorf("reading resource failed, error: %w", err)
	}
	log.Printf("Received quote: %s", resource)

	return nil
}
