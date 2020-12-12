package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

//User is ..
type User struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Dob               string `json:"dob"`
	PhoneNumber       int64  `json:"phonenumber"`
	EmailAddress      string `json:"emailaddress"`
	CreationTimestamp string `json:"creationtimestamp"`
}
type userHandlers struct {
	sync.Mutex
	store map[string]User
}

var globalcounter int64

func (h *userHandlers) usersHandlers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}
func (h *userHandlers) getUser(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h.Lock()

	user, ok := h.store[parts[2]]
	h.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *userHandlers) get(w http.ResponseWriter, r *http.Request) {
	users := make([]User, len(h.store))

	h.Lock()
	i := 0
	for _, user := range h.store {
		users[i] = user
		i++
	}
	h.Unlock()
	jsonBytes, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *userHandlers) post(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
		return
	}

	var user User
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	user.CreationTimestamp = fmt.Sprintf("%s", time.Now().UTC())
	user.ID = fmt.Sprintf("%d", globalcounter)
	globalcounter = globalcounter + 1
	h.Lock()
	h.store[user.ID] = user

	defer h.Unlock()

}

func newUserHandlers() *userHandlers {
	return &userHandlers{
		store: map[string]User{},
	}
}

//Contact Section

//Contact is ..
type Contact struct {
	UserIDOne     string `json:"useridone"`
	UserIDTwo     string `json:"useridtwo"`
	TimeOfContact string `json:"timeofcontact"`
}

type contactHandlers struct {
	sync.Mutex
	store map[string]Contact
}

func (h *contactHandlers) contactHandlers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "POST":
		h.createContact(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

func (h *contactHandlers) createContact(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
		return
	}

	var contact Contact
	err = json.Unmarshal(bodyBytes, &contact)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	contact.TimeOfContact = fmt.Sprintf("%s", time.Now().UTC())
	h.Lock()
	h.store[contact.TimeOfContact] = contact
	defer h.Unlock()
}

func (h *contactHandlers) getContact(w http.ResponseWriter, r *http.Request) {

}

func newContactHandlers() *contactHandlers {
	return &contactHandlers{
		store: map[string]Contact{},
	}
}

func main() {
	globalcounter = 1
	port := ":8080"
	userHandlers := newUserHandlers()
	http.HandleFunc("/users", userHandlers.usersHandlers)
	http.HandleFunc("/users/", userHandlers.getUser)

	//Contact Handle Functions
	contactHandlers := newContactHandlers()
	http.HandleFunc("/contacts", contactHandlers.contactHandlers)
	http.HandleFunc("/contacts/", contactHandlers.getContact)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}
