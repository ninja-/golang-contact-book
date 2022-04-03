package main

import (
	"fmt"
	"log"

	"example.com/contacts/server"
)

func main() {
	fmt.Println("Contacts API server")
	db := server.NewMemoryDatabase()
	if err := db.LoadFixtures(); err != nil {
		log.Fatal(err)
	}

	server, err := server.NewRestServer(db, "./audit.log")
	if err != nil {
		log.Fatal(err)
	}
	server.Start(8080)
}
