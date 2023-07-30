package main

import (
	"log"

	"WordOfWisdom/internal/server"
)

func main() {
	err := server.Run()
	if err != nil {
		log.Fatal(err)
	}
}
