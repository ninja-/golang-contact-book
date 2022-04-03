package server

type ContactDatabase interface {
	// Inserts into the database - in case of id conflict, false will be returned
	Insert(contact Contact) bool
	// Inserts into the database while autogenerating id
	InsertWithNewId(contact Contact) Contact

	// Updates a contact in the database - in case of no matching contact by id, false will be returned
	Update(contact Contact) bool
	// Deletes a contact in the database - in case of no matching contact by id, false will be returned
	Delete(contact Contact) bool

	// Find a contact by id, or returns nil if not found
	FindById(id int) *Contact
	// Find all contacts with last name containg given part. Order is unspecified
	FindByLastNameContains(part string) []Contact
	// Find all contacts matching given email. Order is unspecified
	FindByEmail(email string) []Contact
	// Finds all contacts in the database. Order is unspecified
	FindAll() []Contact
}
