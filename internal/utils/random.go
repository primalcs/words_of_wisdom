package utils

import (
	"encoding/base64"
	"math/rand"
	"strconv"
)

const (
	minRandNumber = 2048
	maxRandNumber = 1024 * 1024
)

func GenerateRandomString() string {
	n := rand.Intn(maxRandNumber) + minRandNumber
	return base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(n)))
}
