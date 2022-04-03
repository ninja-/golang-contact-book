package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type RestServer struct {
	db    ContactDatabase
	audit *log.Logger
}

type auditLog struct {
	Op   string      `json:"op"`
	Data interface{} `json:"data"`
}

func NewRestServer(db ContactDatabase, auditFile string) (*RestServer, error) {
	file, err := os.OpenFile(auditFile, os.O_APPEND|os.O_CREATE, 0700)
	if err != nil {
		return nil, err
	}
	auditLogger := log.New(file, "", log.LstdFlags|log.Lshortfile)
	return &RestServer{db: db, audit: auditLogger}, nil
}

func writeJson(data interface{}, w http.ResponseWriter) error {
	var resp []byte
	resp, err := json.Marshal(data)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "%+v", string(resp))
	return nil
}

type appHandler func(http.ResponseWriter, *http.Request) error

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (r *RestServer) auditLog(op string, data interface{}) {
	log, _ := json.Marshal(auditLog{Op: op, Data: data})
	r.audit.Println(string(log))
}

func (r *RestServer) findAll(w http.ResponseWriter, req *http.Request) error {
	r.auditLog("findAll", nil)
	return writeJson(r.db.FindAll(), w)
}

func (r *RestServer) findById(w http.ResponseWriter, req *http.Request) error {
	idStr := mux.Vars(req)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return errors.New("invalid id")
	}

	r.auditLog("findById", id)

	contact := r.db.FindById(id)
	if contact == nil {
		w.WriteHeader(404)
		return nil
	}

	return writeJson(contact, w)
}

func (r *RestServer) deleteById(w http.ResponseWriter, req *http.Request) error {
	idStr := mux.Vars(req)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return errors.New("invalid id")
	}

	r.auditLog("deleteById", id)

	if !r.db.Delete(Contact{Id: id}) {
		w.WriteHeader(404)
		return nil
	}

	return nil
}

func (r *RestServer) create(w http.ResponseWriter, req *http.Request) error {
	var contact Contact
	if err := json.NewDecoder(req.Body).Decode(&contact); err != nil {
		return err
	}

	if err := contact.Validate(); err != nil {
		return err
	}

	r.auditLog("create", contact.Anonymize())

	contact = r.db.InsertWithNewId(contact)
	return writeJson(contact, w)
}

func (r *RestServer) updateById(w http.ResponseWriter, req *http.Request) error {
	var contact Contact
	if err := json.NewDecoder(req.Body).Decode(&contact); err != nil {
		return err
	}

	if err := contact.Validate(); err != nil {
		return err
	}

	r.auditLog("updateById", contact.Anonymize())

	if !r.db.Update(contact) {
		w.WriteHeader(404)
		return nil
	}

	return nil
}

func (r *RestServer) searchByEmail(w http.ResponseWriter, req *http.Request) error {
	email := mux.Vars(req)["email"]
	r.auditLog("searchByEmail", "*** ANONYMIZED ***")

	contacts := r.db.FindByEmail(email)
	return writeJson(contacts, w)
}

func (r *RestServer) searchByLastNamePart(w http.ResponseWriter, req *http.Request) error {
	lastNamePart := mux.Vars(req)["lastNamePart"]
	r.auditLog("searchByLastNamePart", "*** ANONYMIZED ***")

	contacts := r.db.FindByLastNameContains(lastNamePart)
	return writeJson(contacts, w)
}

func (r *RestServer) Start(port int) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/contacts", appHandler(r.findAll).ServeHTTP).Methods("GET")
	router.HandleFunc("/contacts", appHandler(r.create).ServeHTTP).Methods("POST")
	router.HandleFunc("/contacts/{id}", appHandler(r.findById).ServeHTTP).Methods("GET")
	router.HandleFunc("/contacts/{id}", appHandler(r.deleteById).ServeHTTP).Methods("DELETE")
	router.HandleFunc("/contacts/{id}", appHandler(r.updateById).ServeHTTP).Methods("PUT")

	router.HandleFunc("/contacts/search/email/{email}", appHandler(r.searchByEmail).ServeHTTP).Methods("GET")
	router.HandleFunc("/contacts/search/lastNamePart/{lastNamePart}", appHandler(r.searchByLastNamePart).ServeHTTP).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), router))
}
