package message

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
)

type Type int

const (
	ChallengeRequest Type = iota
	ChallengeResponse
	QuoteResponse
	ErrorResponse
)

var (
	InvalidMessageType = errors.New("invalid message type")
	InvalidPow         = errors.New("invalid PoW")
)

type Message struct {
	Type       Type   `json:"type"`
	Data       string `json:"data"`
	ZerosCount int    `json:"zeros_count"`
}

func NewMessage(t Type, data string, zerosCount int) *Message {
	return &Message{
		Type:       t,
		Data:       data,
		ZerosCount: zerosCount,
	}
}

func (m *Message) SendMessage(conn net.Conn) error {
	v, err := json.Marshal(m)
	if err != nil {
		log.Printf("Error marshalling message: %s\n", err)
		return fmt.Errorf("marshal: %w", err)
	}

	v = append(v, 10)
	_, err = conn.Write(v)
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

func ParseMessage(data []byte) (*Message, error) {
	var msg *Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return msg, nil
}
