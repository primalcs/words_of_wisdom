package server

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	ServerHost = "SERVER_HOST"
	ServerPort = "SERVER_PORT"

	RedisHost = "REDIS_HOST"
	RedisPort = "REDIS_PORT"

	HashCashZerosCount        = "HASHCASH_ZEROS_COUNT"
	HashCashChallengeLifetime = "HASHCASH_CHALLENGE_LIFETIME"
)

type Config struct {
	ServerAddress       string
	RedisAddress        string
	HCZerosCount        int64
	HCChallengeDuration time.Duration
}

func LoadConfig() (*Config, error) {
	var err error
	cfg := &Config{}

	strVal, ok := os.LookupEnv(HashCashZerosCount)
	if !ok {
		return nil, fmt.Errorf("missing %s environment variable", HashCashZerosCount)
	}
	cfg.HCZerosCount, err = strconv.ParseInt(strVal, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid %s environment variable: %s", HashCashZerosCount, strVal)
	}

	strVal, ok = os.LookupEnv(HashCashChallengeLifetime)
	if !ok {
		return nil, fmt.Errorf("missing %s environment variable", HashCashChallengeLifetime)
	}
	cfg.HCChallengeDuration, err = time.ParseDuration(strVal)
	if err != nil {
		return nil, fmt.Errorf("invalid %s environment variable: %s", HashCashChallengeLifetime, strVal)
	}

	strVal, ok = os.LookupEnv(RedisHost)
	if !ok {
		return nil, fmt.Errorf("missing %s environment variable", RedisHost)
	}
	cfg.RedisAddress = strVal + ":"
	strVal, ok = os.LookupEnv(RedisPort)
	if !ok {
		return nil, fmt.Errorf("missing %s environment variable", RedisPort)
	}
	if _, err := strconv.ParseInt(strVal, 10, 64); err != nil {
		return nil, fmt.Errorf("invalid %s environment variable: %s", RedisPort, strVal)
	}
	cfg.RedisAddress += strVal

	strVal, ok = os.LookupEnv(ServerHost)
	if !ok {
		return nil, fmt.Errorf("missing %s environment variable", ServerHost)
	}
	cfg.ServerAddress = strVal + ":"
	strVal, ok = os.LookupEnv(ServerPort)
	if !ok {
		return nil, fmt.Errorf("missing %s environment variable", ServerPort)
	}
	if _, err := strconv.ParseInt(strVal, 10, 64); err != nil {
		return nil, fmt.Errorf("invalid %s environment variable: %s", ServerPort, strVal)
	}
	cfg.ServerAddress += strVal
	return cfg, nil
}
