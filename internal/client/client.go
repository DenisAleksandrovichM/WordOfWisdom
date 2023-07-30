package client

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"WordOfWisdom/pkg/config"
	"WordOfWisdom/pkg/message"
	"WordOfWisdom/pkg/pow"
)

const timeout = 10 * time.Second

func Run() error {
	log.Printf("Connecting to Word of Wisdom server...\n")
	cfg, err := config.Load("config/config.json")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	address := fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort)
	if err = connectAndServe(address); err != nil {
		return fmt.Errorf("run: %w", err)
	}
	return nil
}

func connectAndServe(address string) error {
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}

	return handleConnection(conn)
}

func safeCloseConn(conn net.Conn) {
	if err := conn.Close(); err != nil {
		log.Printf("Error closing conn: %s\n", err)
	}
}

func handleConnection(conn net.Conn) error {
	defer safeCloseConn(conn)
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		requestMsg, err := message.ParseMessage(scanner.Bytes())
		if err != nil {
			return fmt.Errorf("parse message: %w", err)
		}

		switch requestMsg.Type {
		case message.ChallengeRequest:
			if err = handleChallengeRequest(conn, requestMsg); err != nil {
				return fmt.Errorf("handle challenge request: %w", err)
			}
		case message.QuoteResponse:
			return handleQuoteResponse(requestMsg)
		case message.ErrorResponse:
			return handleErrorResponse(requestMsg)
		default:
			return handleDefaultResponse(requestMsg)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from server: %s\n", err)
	}
	return nil
}

func handleChallengeRequest(conn net.Conn, requestMsg *message.Message) error {
	log.Printf("Received challenge request: %s\n", requestMsg.Data)
	challenge := strings.TrimSpace(strings.TrimPrefix(requestMsg.Data, pow.HashCashHeader))
	response, err := pow.SolvePoW(challenge, requestMsg.ZerosCount)
	if err != nil {
		return fmt.Errorf("solve PoW: %w", err)
	}

	msg := message.NewMessage(message.ChallengeResponse, response, 0)
	err = msg.SendMessage(conn)
	if err != nil {
		return fmt.Errorf("send response: %w", err)
	}
	log.Printf("Sent challenge response: %s\n", response)
	return nil
}

func handleQuoteResponse(requestMsg *message.Message) error {
	quote := strings.TrimSpace(requestMsg.Data)
	log.Printf("Received Quote: %s\n", quote)
	return nil
}

func handleErrorResponse(requestMsg *message.Message) error {
	return errors.New(requestMsg.Data)
}

func handleDefaultResponse(_ *message.Message) error {
	return errors.New("unknown message type")
}
