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

	user.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	h.Lock()
	h.store[user.ID] = user
	defer h.Unlock()

}

func newUserHandlers() *userHandlers {
	return &userHandlers{
		store: map[string]User{
			"1": {

				ID:                "1",
				Name:              "Abhishek Soy",
				Dob:               "1997-07-24",
				PhoneNumber:       8603100915,
				EmailAddress:      "soyabhishek81@gmail.com",
				CreationTimestamp: fmt.Sprintf("%s", time.Unix(time.Now().UTC().Unix(), 0)),
			},
		},
	}
}
func main() {
	port := ":8080"
	userHandlers := newUserHandlers()
	http.HandleFunc("/users", userHandlers.usersHandlers)
	http.HandleFunc("/users/", userHandlers.getUser)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}
