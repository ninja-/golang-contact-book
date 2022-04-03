package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"example.com/contacts/server"
)

type HttpClient struct {
	client  *http.Client
	baseUrl string
}

func NewContactsClient(client *http.Client, baseUrl string) *HttpClient {
	return &HttpClient{client: client, baseUrl: baseUrl}
}

func readContact(body io.ReadCloser) (*server.Contact, error) {
	var contact *server.Contact
	defer body.Close()
	if err := json.NewDecoder(body).Decode(&contact); err != nil {
		return nil, err
	}
	return contact, nil
}

func readContactArray(body io.ReadCloser) ([]server.Contact, error) {
	var contacts []server.Contact
	defer body.Close()
	if err := json.NewDecoder(body).Decode(&contacts); err != nil {
		return nil, err
	}
	return contacts, nil
}

func (c *HttpClient) InsertWithNewId(contact server.Contact) (server.Contact, error) {
	body, err := json.Marshal(contact)
	if err != nil {
		return contact, err
	}
	resp, err := c.client.Post(c.baseUrl+"/contacts", "application/json", strings.NewReader(string(body)))
	if err != nil {
		return contact, err
	}

	if resp.StatusCode > 500 {
		return contact, errors.New("unexpected status code " + strconv.Itoa(resp.StatusCode))
	}
	newContact, err := readContact(resp.Body)
	if err != nil {
		return contact, err
	}
	if newContact == nil {
		return contact, errors.New("unexpected nil response")
	}
	return *newContact, nil
}

func (c *HttpClient) Update(contact server.Contact) (bool, error) {
	body, err := json.Marshal(contact)
	if err != nil {
		return false, err
	}
	req, err := http.NewRequest("PUT", c.baseUrl+"/contacts/"+strconv.Itoa(contact.Id), bytes.NewBuffer(body))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return false, err
	}
	if resp.StatusCode > 500 {
		return false, errors.New("unexpected status code " + strconv.Itoa(resp.StatusCode))
	}

	return resp.StatusCode != 404, nil
}

func (c *HttpClient) Delete(contact server.Contact) (bool, error) {
	req, err := http.NewRequest("DELETE", c.baseUrl+"/contacts/"+strconv.Itoa(contact.Id), nil)
	if err != nil {
		return false, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return false, err
	}
	if resp.StatusCode > 500 {
		return false, errors.New("unexpected status code " + strconv.Itoa(resp.StatusCode))
	}

	return resp.StatusCode != 404, nil
}

func (c *HttpClient) FindById(id int) (*server.Contact, error) {
	resp, err := c.client.Get(c.baseUrl + "/contacts/" + strconv.Itoa(id))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 500 {
		return nil, errors.New("unexpected status code " + strconv.Itoa(resp.StatusCode))
	}
	if resp.StatusCode == 404 {
		return nil, nil
	}

	return readContact(resp.Body)
}

func (c *HttpClient) FindByLastNameContains(lastNamePart string) ([]server.Contact, error) {
	resp, err := c.client.Get(c.baseUrl + "/contacts/search/lastNamePart/" + url.QueryEscape(lastNamePart))
	if err != nil {
		return nil, err
	}
	return readContactArray(resp.Body)
}

func (c *HttpClient) FindByEmail(email string) ([]server.Contact, error) {
	resp, err := c.client.Get(c.baseUrl + "/contacts/search/email/" + url.QueryEscape(email))
	if err != nil {
		return nil, err
	}
	return readContactArray(resp.Body)
}

func (c *HttpClient) FindAll() ([]server.Contact, error) {
	resp, err := c.client.Get(c.baseUrl + "/contacts")
	if err != nil {
		return nil, err
	}
	return readContactArray(resp.Body)
}

var _ Client = (*HttpClient)(nil)
