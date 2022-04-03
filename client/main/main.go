package main

import (
	"net/http"
	"os"

	"example.com/contacts/client"
)

func main() {
	endpoint := "http://localhost:8080"
	restClient := client.NewContactsClient(&http.Client{}, endpoint)
	cli := client.NewCliClient(restClient)
	cli.HandleCommand(os.Args[1:])
}
