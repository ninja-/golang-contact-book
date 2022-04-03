package server

import (
	"strings"

	"gopkg.in/yaml.v2"
)

type MemoryDatabase struct {
	data      map[int]Contact
	highestId int
}

var _ ContactDatabase = (*MemoryDatabase)(nil)

func NewMemoryDatabase() *MemoryDatabase {
	data := make(map[int]Contact)
	return &MemoryDatabase{
		data: data,
	}
}

func (m *MemoryDatabase) hasContact(id int) bool {
	_, ok := m.data[id]
	return ok
}

func (m *MemoryDatabase) dataCopy() []Contact {
	// Order is unspecified
	var result []Contact
	for _, contact := range m.data {
		result = append(result, contact)
	}
	return result
}

func (m *MemoryDatabase) Insert(contact Contact) bool {
	if m.hasContact(contact.Id) {
		return false
	}
	if contact.Id > m.highestId {
		m.highestId = contact.Id
	}

	m.data[contact.Id] = contact
	return true
}

func (m *MemoryDatabase) InsertWithNewId(contact Contact) Contact {
	contact.Id = m.highestId + 1
	m.Insert(contact)
	return contact
}

func (m *MemoryDatabase) Delete(contact Contact) bool {
	if !m.hasContact(contact.Id) {
		return false
	}

	delete(m.data, contact.Id)
	return true
}

func (m *MemoryDatabase) Update(contact Contact) bool {
	if !m.hasContact(contact.Id) {
		return false
	}

	m.data[contact.Id] = contact
	return true
}

func (m *MemoryDatabase) FindAll() []Contact {
	return m.dataCopy()
}

func (m *MemoryDatabase) FindById(id int) *Contact {
	if !m.hasContact(id) {
		return nil
	}

	contact := m.data[id]
	return contact.Clone()
}

func (m *MemoryDatabase) FindByEmail(email string) []Contact {
	var result []Contact

	for _, contact := range m.data {
		if contact.Email == email {
			result = append(result, contact)
		}
	}

	return result
}

func (m *MemoryDatabase) FindByLastNameContains(part string) []Contact {
	var result []Contact

	for _, contact := range m.data {
		if strings.Contains(contact.LastName, part) {
			result = append(result, contact)
		}
	}

	return result
}

var fixtures = `
- name: John
  lastName: Lennon
  email: john.lennon@thebeatles.com
  id: 1
- name: Paul
  lastName : McCartney
  email: paul.mccartney@thebeatles.com
  id: 2
- name: George
  lastName: Harrison
  email: georger.harrison@thebeatles.com
  id: 3
- name: Ringo
  lastName: Starr
  email: ringo.starr@thebeatles.com
  id: 4`

func (m *MemoryDatabase) LoadFixtures() error {
	var data []Contact

	err := yaml.Unmarshal([]byte(fixtures), &data)
	if err != nil {
		return err
	}

	for _, contact := range data {
		m.Insert(contact)
	}

	return nil
}
