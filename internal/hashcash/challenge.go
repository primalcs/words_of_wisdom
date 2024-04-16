package hashcash

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/primalcs/words_of_wisdom/pkg/interfaces"
)

var _ interfaces.POWChallenge = &Challenge{}

type Challenge struct {
	Version    int64
	ZerosCount int64
	Date       time.Time
	Resource   string
	Extension  string
	Rand       string
	Counter    int64
}

func (h *Challenge) GetDate() time.Time {
	return h.Date
}

func (h *Challenge) GetRand() string {
	return h.Rand
}

func (h *Challenge) GetResource() string {
	return h.Resource
}

func (h *Challenge) encode() string {
	if h == nil {
		return ""
	}
	stringDate := h.Date.Format("060102")
	return fmt.Sprintf("%d:%d:%s:%s:%s:%s:%d", h.Version, h.ZerosCount, stringDate, h.Resource, h.Extension, h.Rand, h.Counter)
}

func countSHA1(data string) (string, error) {
	h := sha1.New()
	_, err := h.Write([]byte(data))
	if err != nil {
		return "", fmt.Errorf("could not write data, error: %w", err)
	}
	return string(h.Sum(nil)), nil
}

func isHashCorrect(hash string, zerosCount int64) bool {
	if zerosCount > int64(len(hash)) {
		return false
	}
	for _, ch := range hash[:zerosCount] {
		if ch != 0x30 {
			return false
		}
	}
	return true
}

func (h *Challenge) Solve(maxIterationsAmount int64) error {
	for h.Counter <= maxIterationsAmount {
		strHeader := h.encode()
		hash, err := countSHA1(strHeader)
		if err != nil {
			return fmt.Errorf("could not hash header, error: %w", err)
		}
		if isHashCorrect(hash, h.ZerosCount) {
			return nil
		}
		h.Counter++
	}
	return errors.New("max iterations exceeded")
}

func (h *Challenge) Verify() (bool, error) {
	hash, err := countSHA1(h.encode())
	if err != nil {
		return false, fmt.Errorf("computing sha failed: %v", err)
	}
	return isHashCorrect(hash, h.ZerosCount), nil
}

func (h *Challenge) ToJSON() (string, error) {
	headerBytes, err := json.Marshal(h)
	if err != nil {
		return "", fmt.Errorf("hashcash marshalling failed: %w", err)
	}
	return string(headerBytes), nil
}

func (h *Challenge) FromJSON(jsonString string) error {
	return json.Unmarshal([]byte(jsonString), &h)
}
