package client

import "example.com/contacts/server"

type Client interface {
	InsertWithNewId(contact server.Contact) (server.Contact, error)

	Update(contact server.Contact) (bool, error)
	Delete(contact server.Contact) (bool, error)

	FindById(id int) (*server.Contact, error)
	FindByLastNameContains(part string) ([]server.Contact, error)
	FindByEmail(email string) ([]server.Contact, error)
	FindAll() ([]server.Contact, error)
}
