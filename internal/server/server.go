package server

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/primalcs/words_of_wisdom/internal/storage/book"

	"github.com/primalcs/words_of_wisdom/internal/hashcash"
	"github.com/primalcs/words_of_wisdom/internal/interfaces"
	"github.com/primalcs/words_of_wisdom/internal/storage/cache"
	pkg "github.com/primalcs/words_of_wisdom/pkg/interfaces"
	"github.com/primalcs/words_of_wisdom/pkg/protocol"
)

type Server struct {
	listener net.Listener
	cache    interfaces.Cache
	builder  pkg.POWChallengeBuilder
	book     interfaces.Book
}

func NewServer(ctx context.Context, cfg *Config) (*Server, error) {
	listener, err := net.Listen("tcp", cfg.ServerAddress)
	if err != nil {
		return nil, err
	}

	redisClient, err := cache.NewRedisCache(ctx, cfg.RedisAddress)
	if err != nil {
		return nil, fmt.Errorf("error init cache, error: %w", err)
	}

	ws, err := book.NewWisdom()
	if err != nil {
		return nil, fmt.Errorf("error init cache, error: %w", err)
	}

	return &Server{
		cache:    redisClient,
		listener: listener,
		book:     ws,

		builder: &hashcash.Builder{
			ZerosCount:        cfg.HCZerosCount,
			ChallengeDuration: cfg.HCChallengeDuration,
		},
	}, nil
}

func (srv *Server) Close(ctx context.Context) error {
	if err := srv.cache.Close(ctx); err != nil {
		return err
	}
	if err := srv.listener.Close(); err != nil {
		return err
	}
	return nil
}

func (srv *Server) Run(ctx context.Context) error {
	for {
		conn, err := srv.listener.Accept()
		if err != nil {
			return fmt.Errorf("error accepting connection: %w", err)
		}
		go srv.handleConnection(ctx, conn)
	}
}

func (srv *Server) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	log.Printf("handling connection: %s", conn.RemoteAddr().String())
	reader := bufio.NewReader(conn)
	for {
		msg, err := protocol.Read(reader)
		if err != nil {
			log.Printf("reading message from connection failed: %v", err)
			return
		}
		payload := srv.book.GetOne()
		if err = srv.processRequest(ctx, conn, conn.RemoteAddr().String(), msg, payload); err != nil {
			log.Printf("process request error: %v", err)
			return
		}
	}
}

func (srv *Server) sendChallengeToClient(ctx context.Context, connWriter io.Writer, clientId string) error {
	currentTime := time.Now()
	connBufWriter := bufio.NewWriter(connWriter)

	challenge, randVal := srv.builder.GenerateRandomChallenge(currentTime, clientId)
	// to check in the future that server actually used this randVal
	err := srv.cache.InsertClientToken(ctx, clientId, randVal, srv.builder.GetChallengeDuration())
	if err != nil {
		return fmt.Errorf("inserting into storage set failed, error: %w", err)
	}

	return protocol.SendChallenge(connBufWriter, challenge)
}

func (srv *Server) checkClientChallengeSolution(ctx context.Context, clientId string, solution pkg.POWChallenge) (bool, error) {
	if solution.GetResource() != clientId {
		return false, errors.New("invalid resource")
	}

	clientToken, err := srv.cache.GetClientToken(ctx, solution.GetResource())
	if err != nil {
		return false, fmt.Errorf("client token doesn't exist, error: %w", err)
	}
	if clientToken != solution.GetRand() {
		return false, errors.New("wrong client token")
	}

	if err := srv.cache.Delete(ctx, solution.GetResource()); err != nil {
		return false, fmt.Errorf("deleting from storage set failed, error: %w", err)
	}

	// solution shouldn't take more than setup challengeDuration
	if time.Now().Sub(solution.GetDate()) > srv.builder.GetChallengeDuration() {
		return false, errors.New("challenge expired")
	}

	// verify according to POW puzzle
	valid, err := solution.Verify()
	if err != nil {
		return false, fmt.Errorf("verifying challenge solution failed, error: %w", err)
	}
	return valid, err
}

func (srv *Server) processRequest(ctx context.Context, connWriter io.Writer, clientId string,
	msg *protocol.Message, payload string) error {
	switch msg.MessageType {
	case protocol.ForceQuit:
		log.Print("client requests to close connection")
		return nil

	case protocol.ChallengeRequest:
		log.Printf("client %s requests challenge\n", clientId)
		return srv.sendChallengeToClient(ctx, connWriter, clientId)
	case protocol.ResourceRequest:
		powSolution := &hashcash.Challenge{}
		if err := powSolution.FromJSON(msg.Payload); err != nil {
			return err
		}

		log.Printf("client solved challenge and requests resource. client: %s, payload %s\n", clientId, msg.Payload)

		valid, err := srv.checkClientChallengeSolution(ctx, clientId, powSolution)
		if err != nil {
			return fmt.Errorf("verifying challenge solution failed, error: %w", err)
		}
		if !valid {
			return errors.New("invalid solution")
		}

		log.Print("Solution verified. Sending a word of wisdom.")

		return protocol.SendResource(bufio.NewWriter(connWriter), payload)
	default:
		return errors.New("invalid message type")
	}
}
