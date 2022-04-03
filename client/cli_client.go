package client

import (
	"log"
	"strconv"

	"example.com/contacts/server"
)

type CliClient struct {
	client Client
}

func NewCliClient(client Client) *CliClient {
	return &CliClient{
		client: client,
	}
}

func (c *CliClient) HandleCommand(args []string) {
	if len(args) < 1 {
		log.Print("Usage: ./client <add|delete|update|findByEmail|findByLastNamePart> [...]")
		return
	}

	// Create, delete, update
	// Search by email, search by lastName part

	if args[0] == "add" {
		if len(args) < 4 {
			log.Print("Usage: ./client add <name> <lastName> <email>")
			return
		}

		name, lastName, email := args[1], args[2], args[3]
		_, err := c.client.InsertWithNewId(server.Contact{
			Name:     name,
			LastName: lastName,
			Email:    email,
		})
		if err != nil {
			log.Print(err)
			return
		}
		log.Print("Sucessfully created contact")
		return
	}

	if args[0] == "delete" {
		if len(args) < 2 {
			log.Print("Usage: ./client delete <id>")
			return
		}

		id, err := strconv.Atoi(args[1])
		if err != nil {
			log.Print(err)
			return
		}

		resp, err := c.client.Delete(server.Contact{Id: id})
		if err != nil {
			log.Print(err)
			return
		}

		if !resp {
			log.Print("Failed to delete contact - not found by id")
		} else {
			log.Print("Successfully deleted contact")
		}
		return
	}

	if args[0] == "update" {
		if len(args) < 5 {
			log.Print("Usage: ./client update <id> <name> <lastName> <email>")
			return
		}

		id, err := strconv.Atoi(args[1])
		if err != nil {
			log.Print(err)
			return
		}

		name, lastName, email := args[2], args[3], args[4]

		resp, err := c.client.Update(server.Contact{
			Id:       id,
			Name:     name,
			LastName: lastName,
			Email:    email,
		})

		if err != nil {
			log.Print(err)
			return
		}

		if !resp {
			log.Print("Failed to update contact - not found by id")
		} else {
			log.Print("Successfully updated contact")
		}
		return
	}

	if args[0] == "findByEmail" {
		if len(args) < 1 {
			log.Print("Usage: ./client findByEmail <email>")
			return
		}
		email := args[1]
		resp, err := c.client.FindByEmail(email)

		if err != nil {
			log.Print(err)
			return
		}

		if resp == nil {
			log.Print("Contacts not found")
			return
		}

		log.Print("Found contacts: ", resp)
		return
	}

	if args[0] == "findByLastNamePart" {
		if len(args) < 1 {
			log.Print("Usage: ./client findByLastNamePart <part>")
			return
		}
		part := args[1]
		resp, err := c.client.FindByLastNameContains(part)

		if err != nil {
			log.Print(err)
			return
		}
		if resp == nil {
			log.Print("Contact not found")
			return
		}
		log.Print("Found contacts: ", resp)
		return
	}

	log.Print("Unknown command")
}
