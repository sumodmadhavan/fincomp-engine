package main

import (
	"log"

	"financialapi/internal/api"
)

func main() {
	server := api.NewServer()
	log.Fatal(server.Run(":8080"))
}
