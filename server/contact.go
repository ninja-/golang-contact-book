package server

import "errors"

type Contact struct {
	Id int `json:"id" yaml:"id"`

	Name     string `json:"name" yaml:"name"`
	LastName string `json:"lastName" yaml:"lastName"`
	Email    string `json:"email" yaml:"email"`
}

func (c *Contact) Clone() *Contact {
	return &Contact{
		Id: c.Id,

		Name:     c.Name,
		LastName: c.LastName,
		Email:    c.Email,
	}
}

func (c *Contact) Validate() error {
	// Or could use validator library in future

	if c.Name == "" {
		return errors.New("empty contact name")
	}
	if c.LastName == "" {
		return errors.New("empty contact last name")
	}
	if c.Email == "" {
		return errors.New("empty contact email")
	}

	return nil
}

func (c *Contact) Anonymize() Contact {
	return Contact{
		Id:       c.Id,
		Email:    "*** ANONYMIZED ***",
		Name:     "*** ANONYMIZED ***",
		LastName: "*** ANONYMIZED ***",
	}
}
