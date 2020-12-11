package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type User struct {
	Name string `json:"name"`
	ID   string `json:"id"`
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
		store: map[string]User{},
	}
}
func main() {
	port := ":8080"
	userHandlers := newUserHandlers()
	http.HandleFunc("/users", userHandlers.usersHandlers)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}
