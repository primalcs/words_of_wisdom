package protocol

import (
	"bufio"
	"fmt"

	"github.com/primalcs/words_of_wisdom/pkg/interfaces"
)

func RequestChallenge(bufWriter *bufio.Writer) error {
	return Write(bufWriter, &Message{
		MessageType: ChallengeRequest,
	})
}

func SendChallenge(bufWriter *bufio.Writer, puzzle interfaces.POWChallenge) error {
	payload, err := puzzle.ToJSON()
	if err != nil {
		return fmt.Errorf("marshalling puzzle failed: %v", err)
	}
	msg := &Message{
		MessageType: ChallengeResponse,
		Payload:     payload,
	}
	return Write(bufWriter, msg)
}

func ReceiveChallenge(bufReader *bufio.Reader, builder interfaces.POWChallenge) (interfaces.POWChallenge, error) {
	msg, err := Read(bufReader)
	if err != nil {
		return nil, fmt.Errorf("reading challenge msg failed: %w", err)
	}
	if err = builder.FromJSON(msg.Payload); err != nil {
		return nil, fmt.Errorf("hashcash unmarshal failed: %w", err)
	}
	return builder, nil
}

func SendChallengeSolution(bufWriter *bufio.Writer, solution interfaces.POWChallenge) error {
	payload, err := solution.ToJSON()
	if err != nil {
		return err
	}
	err = Write(bufWriter, &Message{
		MessageType: ResourceRequest,
		Payload:     payload,
	})
	if err != nil {
		return fmt.Errorf("sending solution failed: %v", err)
	}
	return nil
}

func SendResource(bufWriter *bufio.Writer, payload string) error {
	msg := &Message{
		MessageType: ResourceResponse,
		Payload:     payload,
	}

	return Write(bufWriter, msg)
}

func ReceiveResource(bufReader *bufio.Reader) (string, error) {
	msgWithResource, err := Read(bufReader)
	if err != nil {
		return "", fmt.Errorf("receiving resource failed: %v", err)
	}
	return msgWithResource.Payload, nil
}
