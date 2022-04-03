package client

import (
	"example.com/contacts/server"
	"github.com/stretchr/testify/mock"
)

type ClientMock struct {
	mock.Mock
}

func (c *ClientMock) InsertWithNewId(contact server.Contact) (server.Contact, error) {
	args := c.Called(contact)
	return args.Get(0).(server.Contact), args.Error(1)
}

func (c *ClientMock) Update(contact server.Contact) (bool, error) {
	args := c.Called(contact)
	return args.Bool(0), args.Error(1)
}

func (c *ClientMock) Delete(contact server.Contact) (bool, error) {
	args := c.Called(contact)
	return args.Bool(0), args.Error(1)
}

func (c *ClientMock) FindById(id int) (*server.Contact, error) {
	args := c.Called(id)
	return args.Get(0).(*server.Contact), args.Error(1)
}

func (c *ClientMock) FindByLastNameContains(lastNamePart string) ([]server.Contact, error) {
	args := c.Called(lastNamePart)
	return args.Get(0).([]server.Contact), args.Error(1)
}

func (c *ClientMock) FindByEmail(email string) ([]server.Contact, error) {
	args := c.Called(email)
	return args.Get(0).([]server.Contact), args.Error(1)
}

func (c *ClientMock) FindAll() ([]server.Contact, error) {
	args := c.Called()
	return args.Get(0).([]server.Contact), args.Error(1)
}

var _ Client = (*ClientMock)(nil)
