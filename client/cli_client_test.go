package client_test

import (
	"testing"

	"example.com/contacts/client"
	"example.com/contacts/server"
)

func TestAddContact(t *testing.T) {
	mock := &client.ClientMock{}
	cli := client.NewCliClient(mock)

	contact := server.Contact{
		Name:     "name",
		LastName: "lastName",
		Email:    "email@email.com",
	}

	mock.On("InsertWithNewId", contact).Return(contact, nil)

	cli.HandleCommand([]string{"add", contact.Name, contact.LastName, contact.Email})
	mock.AssertExpectations(t)
}

func TestDeleteContact(t *testing.T) {
	mock := &client.ClientMock{}
	cli := client.NewCliClient(mock)

	contact := server.Contact{
		Id: 1,
	}

	mock.On("Delete", contact).Return(true, nil)

	cli.HandleCommand([]string{"delete", "1", contact.Name, contact.LastName, contact.Email})
	mock.AssertExpectations(t)
}

func TestUpdateContact(t *testing.T) {
	mock := &client.ClientMock{}
	cli := client.NewCliClient(mock)

	contact := server.Contact{
		Id:       1,
		Name:     "name",
		LastName: "lastName",
		Email:    "email@email.com",
	}

	mock.On("Update", contact).Return(true, nil)

	cli.HandleCommand([]string{"update", "1", contact.Name, contact.LastName, contact.Email})
	mock.AssertExpectations(t)
}

func TestFindByEmail(t *testing.T) {
	mock := &client.ClientMock{}
	cli := client.NewCliClient(mock)

	contact := server.Contact{
		Id:       1,
		Name:     "name",
		LastName: "lastName",
		Email:    "email@email.com",
	}

	mock.On("FindByEmail", contact.Email).Return([]server.Contact{contact}, nil)

	cli.HandleCommand([]string{"findByEmail", contact.Email})
	mock.AssertExpectations(t)
}

func TestFindByLastNamePart(t *testing.T) {
	mock := &client.ClientMock{}
	cli := client.NewCliClient(mock)

	contact := server.Contact{
		Id:       1,
		Name:     "name",
		LastName: "lastName",
		Email:    "email@email.com",
	}

	mock.On("FindByLastNameContains", contact.LastName).Return([]server.Contact{contact}, nil)

	cli.HandleCommand([]string{"findByLastNamePart", contact.LastName})
	mock.AssertExpectations(t)
}
