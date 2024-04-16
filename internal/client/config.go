package client

import (
	"fmt"
	"os"
	"strconv"
)

const (
	HashcashMaxIterations = "HASHCASH_MAX_ITERATIONS"
	ServerHost            = "SERVER_HOST"
	ServerPort            = "SERVER_PORT"
)

type Config struct {
	ServerAddress string
	MaxIterations int64
}

func LoadConfig() (*Config, error) {
	var err error
	cfg := &Config{}
	it, ok := os.LookupEnv(HashcashMaxIterations)
	if !ok {
		return nil, fmt.Errorf("missing %s environment variable", HashcashMaxIterations)
	}
	cfg.MaxIterations, err = strconv.ParseInt(it, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid %s environment variable: %s", HashcashMaxIterations, err)
	}
	pt, ok := os.LookupEnv(ServerPort)
	if !ok {
		return nil, fmt.Errorf("missing %s environment variable", ServerPort)
	}
	if _, err = strconv.Atoi(pt); err != nil {
		return nil, fmt.Errorf("invalid %s environment variable: %s", ServerPort, err)
	}
	cfg.ServerAddress = os.Getenv(ServerHost) + ":" + pt
	return cfg, nil
}
