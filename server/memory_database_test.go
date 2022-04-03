package server_test

import (
	"sort"
	"testing"

	"example.com/contacts/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sortContactsById(contacts []server.Contact) []server.Contact {
	sort.Slice(contacts, func(i, j int) bool {
		return contacts[i].Id < contacts[j].Id
	})
	return contacts
}

func createDatabaset(t *testing.T) *server.MemoryDatabase {
	db := server.NewMemoryDatabase()
	return db
}

func TestInsertNormalAndConflict(t *testing.T) {
	db := createDatabaset(t)
	contact := server.Contact{
		Id:       1,
		Name:     "Test",
		LastName: "test",
		Email:    "test@test.com",
	}

	ret := db.Insert(contact)
	assert.True(t, ret, "wrong return value without conflict")
	ret = db.Insert(contact)
	assert.False(t, ret, "wrong return value with conflict")

	assert.Equal(t, &contact, db.FindById(contact.Id))
}

func TestInsertWithId(t *testing.T) {
	db := createDatabaset(t)
	contact := server.Contact{
		Name:     "Test",
		LastName: "test",
		Email:    "test@test.com",
	}

	ret := db.InsertWithNewId(contact)
	assert.Equal(t, 1, ret.Id)
	ret = db.InsertWithNewId(contact)
	assert.Equal(t, 2, ret.Id)
}

func TestUpdateNormal(t *testing.T) {
	db := createDatabaset(t)
	contact := server.Contact{
		Id:       1,
		Name:     "Test",
		LastName: "test",
		Email:    "test@test.com",
	}

	ret := db.Insert(contact)
	require.True(t, ret)

	contact.Name = "Test2"
	ret = db.Update(contact)
	require.True(t, ret)

	newContact := db.FindById(contact.Id)
	require.NotNil(t, newContact)

	assert.Equal(t, "Test2", newContact.Name)
	assert.Equal(t, &contact, newContact)
}

func TestUpdateNoMatch(t *testing.T) {
	db := createDatabaset(t)
	contact := server.Contact{
		Id:       1,
		Name:     "Test",
		LastName: "test",
		Email:    "test@test.com",
	}
	ret := db.Update(contact)
	assert.False(t, ret)
}

func TestDeleteNormalAndNoMatch(t *testing.T) {
	db := createDatabaset(t)
	contact := server.Contact{
		Id:       1,
		Name:     "Test",
		LastName: "test",
		Email:    "test@test.com",
	}

	ret := db.Insert(contact)
	require.True(t, ret)

	ret = db.Delete(contact)
	assert.True(t, ret, "wrong return value with match")

	ret = db.Delete(contact)
	assert.False(t, ret, "wrong return value without match")
}

func TestObjectChangeDoesNotAffectDatabase(t *testing.T) {
	db := createDatabaset(t)
	contact := server.Contact{
		Id:       1,
		Name:     "Test",
		LastName: "test",
		Email:    "test@test.com",
	}

	db.Insert(contact)
	contact.Name = "Test2"
	ret := db.FindById(contact.Id)
	ret.Name = "Test2"
	ret = db.FindById(contact.Id)

	assert.Equal(t, "Test", ret.Name)
}

func TestFindAll(t *testing.T) {
	db := createDatabaset(t)
	contacts := []server.Contact{
		{
			Id:       1,
			Name:     "Test",
			LastName: "test",
			Email:    "test@test.com",
		},
		{
			Id:       2,
			Name:     "Test2",
			LastName: "test",
			Email:    "test2@test.com",
		},
	}

	db.Insert(contacts[0])
	db.Insert(contacts[1])

	// Order not guaranteed
	ret := db.FindAll()
	ret = sortContactsById(ret)

	assert.Equal(t, contacts, ret)
}

func TestFindById(t *testing.T) {
	db := createDatabaset(t)
	contacts := []server.Contact{
		{
			Id:       1,
			Name:     "Test",
			LastName: "test",
			Email:    "test@test.com",
		},
		{
			Id:       2,
			Name:     "Test2",
			LastName: "test",
			Email:    "test2@test.com",
		},
	}

	db.Insert(contacts[0])
	db.Insert(contacts[1])

	assert.Equal(t, &contacts[1], db.FindById(contacts[1].Id))
}

func TestFindByEmailMatchAndNoMatch(t *testing.T) {
	db := createDatabaset(t)
	contact := server.Contact{
		Id:       1,
		Name:     "Test",
		LastName: "test",
		Email:    "test@test.com",
	}

	db.Insert(contact)

	ret := db.FindByEmail("no_match")
	assert.Equal(t, 0, len(ret), "wrong value without match")

	ret = db.FindByEmail("test@test.com")
	assert.Equal(t, []server.Contact{contact}, ret, "wrong value with match")
}

func TestFindByLastNameContains(t *testing.T) {
	db := createDatabaset(t)
	contacts := []server.Contact{
		{
			Id:       1,
			Name:     "Test",
			LastName: "test_SUBSTR_test",
			Email:    "test@test.com",
		},
		{
			Id:       2,
			Name:     "Test2",
			LastName: "test_SUBSTR_test",
			Email:    "test2@test.com",
		},
		{
			Id:       3,
			Name:     "Test3",
			LastName: "test",
			Email:    "test3@test.com",
		},
	}

	db.Insert(contacts[0])
	db.Insert(contacts[1])
	db.Insert(contacts[2])

	// Order not guaranteed
	ret := db.FindByLastNameContains("SUBSTR")

	ret = sortContactsById(ret)

	require.Equal(t, 2, len(ret))
	assert.Equal(t, contacts[0], ret[0])
	assert.Equal(t, contacts[1], ret[1])
}
