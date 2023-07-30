package server

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"

	"WordOfWisdom/internal/server/quotes"
	"WordOfWisdom/pkg/config"
	"WordOfWisdom/pkg/message"
	"WordOfWisdom/pkg/pow"
)

func Run() error {
	log.Println("Starting Word of Wisdom server...")
	cfg, err := config.Load("config/config.json")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	address := fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	defer safeCloseListener(listener)

	counter := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting conn: %s\n", err)
			continue
		}

		counter++
		zerosCount := cfg.HashcashZerosCount + counter/cfg.IncZerosCountLimit
		go connectionHandler(conn, zerosCount)
		counter--
	}
}

func safeCloseListener(listener net.Listener) {
	if err := listener.Close(); err != nil {
		log.Printf("Error closing listener: %s\n", err)
	}
}

func connectionHandler(conn net.Conn, zerosCount int) {
	defer safeCloseConn(conn)
	log.Printf("New conn: %s\n", conn.RemoteAddr())
	challengeMsg := createChallengeMsg(zerosCount)
	if err := challengeMsg.SendMessage(conn); err != nil {
		log.Printf("Error sending message: %s\n", err)
		return
	}

	handleConnectionResponse(conn, challengeMsg)
}

func safeCloseConn(conn net.Conn) {
	log.Printf("Closing conn: %s\n", conn.RemoteAddr())
	if err := conn.Close(); err != nil {
		log.Printf("Error closing conn: %s\n", err)
	}
}

func createChallengeMsg(zerosCount int) *message.Message {
	return message.NewMessage(message.ChallengeRequest, fmt.Sprintf("%s%d", pow.HashCashHeader, rand.Intn(1000)), zerosCount)
}

func handleConnectionResponse(conn net.Conn, challengeMsg *message.Message) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		if requestMsg, err := message.ParseMessage(scanner.Bytes()); err != nil {
			log.Printf("Error parsing response: %s\n", err)
			return
		} else {
			handleRequest(conn, requestMsg, challengeMsg)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from client: %s\n", err)
	}
}

func handleRequest(conn net.Conn, requestMsg, challengeMsg *message.Message) {
	var msg *message.Message
	if requestMsg.Type == message.ChallengeResponse {
		log.Printf("Received challenge response from %s: %s\n", conn.RemoteAddr(), requestMsg.Data)
		msg = handleChallengeResponse(requestMsg, challengeMsg)
	} else {
		log.Printf("Received invalid response from %s: %s\n", conn.RemoteAddr(), requestMsg.Data)
		msg = handleInvalidResponse()
	}

	if err := msg.SendMessage(conn); err != nil {
		log.Printf("Error sending response to %s: %s\n", conn.RemoteAddr(), err)
		return
	}
}

func handleChallengeResponse(requestMsg, challengeMsg *message.Message) *message.Message {
	hash := pow.CalculateHash(challengeMsg.Data, requestMsg.Data)
	if pow.IsCorrect(hash, requestMsg.ZerosCount) {
		return message.NewMessage(message.QuoteResponse, quotes.GetRandomQuote(), 0)
	}

	return message.NewMessage(message.ErrorResponse, message.InvalidPow.Error(), 0)
}

func handleInvalidResponse() *message.Message {
	return message.NewMessage(message.ErrorResponse, message.InvalidMessageType.Error(), 0)
}
