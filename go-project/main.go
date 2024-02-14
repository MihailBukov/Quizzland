package main

import (
	"log"
)

func main() {
	err := NewMySqlStore()
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	server := NewApiServer(":3000")
	server.Run()
}
