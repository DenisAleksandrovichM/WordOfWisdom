package main

import (
	"log"

	"WordOfWisdom/internal/client"
)

func main() {
	err := client.Run()
	if err != nil {
		log.Fatal(err)
	}
}
