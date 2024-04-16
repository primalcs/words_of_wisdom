package protocol

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
)

type MessageType uint8

const (
	ForceQuit MessageType = iota + 1
	ChallengeRequest
	ChallengeResponse
	ResourceRequest
	ResourceResponse
)

type Message struct {
	MessageType MessageType `json:"type"`
	Payload     string      `json:"payload,omitempty"`
}

func Read(reader *bufio.Reader) (*Message, error) {
	ba, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("could not read message, error: %w", err)
	}

	var msg Message
	err = json.Unmarshal(ba, &msg)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshall message, error: %w", err)
	}
	return &msg, nil
}

func Write(writer *bufio.Writer, msg *Message) error {
	if msg == nil {
		return errors.New("message is nil")
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("could not marshall message, error: %w", err)
	}
	msgBytes = append(msgBytes, '\n')
	if _, err := writer.Write(msgBytes); err != nil {
		return fmt.Errorf("could not write response, error: %w", err)
	}
	return writer.Flush()
}
