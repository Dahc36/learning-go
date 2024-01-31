package main

import (
	"errors"
	"fmt"
	"net/http"
)

// -- Controller --
type Logger interface {
	Log(message string)
}

type Logic interface {
	SayHello(userId string) (string, error)
}

type Controller struct {
	l     Logger
	logic Logic
}

func (c Controller) Greet(w http.ResponseWriter, r *http.Request) {
	c.l.Log("In Greet")
	userId := r.URL.Query().Get("user_id")
	message, err := c.logic.SayHello(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(message))
}

func NewController(l Logger, logic Logic) Controller {
	return Controller{l: l, logic: logic}
} // -- End - Controller --

// My types
type UserId = string
type User struct {
	Id        int
	FirstName string
	LastName  string
}

type DataStore interface {
	GetUser(userId UserId) (User, bool)
}

// My SimpleLogic
type SimpleLogic struct {
	l  Logger
	ds DataStore
}

// SimpleLogic implements Logic implicitly
func (sl SimpleLogic) SayHello(userId UserId) (string, error) {
	sl.l.Log("in sl SayHello for " + userId)
	user, ok := sl.ds.GetUser(userId)
	if !ok {
		return "", errors.New("unknown user")
	}

	return fmt.Sprintf("Hello %v %v\n", user.FirstName, user.LastName), nil
}

func NewSimpleLogic(l Logger, ds DataStore) SimpleLogic {
	return SimpleLogic{l: l, ds: ds}
}

// My SimpleDataStore
type SimpleDataStore struct {
	userData map[UserId]User
}

// Implementing DataStore implicitly
func (sd SimpleDataStore) GetUser(userId UserId) (User, bool) {
	user, ok := sd.userData[userId]
	return user, ok
}

func NewSimpleDataStore() SimpleDataStore {
	return SimpleDataStore{
		userData: map[UserId]User{
			"1": {Id: 1, FirstName: "Fred", LastName: "Astaire"},
			"2": {Id: 2, FirstName: "Ginger", LastName: "Rogers"},
			"3": {Id: 3, FirstName: "Audrey", LastName: "Hepburn"},
		},
	}
}

// My Logger
func myLogOutput(message string) {
	fmt.Println(message)
}

// LoggerAdapter is of the same underlying type as myLogOutput
// That enables the type conversion l := LoggerAdapter(myLogOutput)
type LoggerAdapter func(message string)

// LoggerAdapter implements Logger implicitly
func (lg LoggerAdapter) Log(message string) {
	lg(message)
}

// Main
func main() {
	l := LoggerAdapter(myLogOutput)
	ds := NewSimpleDataStore()
	logic := NewSimpleLogic(l, ds)
	c := NewController(l, logic)
	http.HandleFunc("/hello", c.Greet)
	fmt.Println("Listening to port 8000...")
	http.ListenAndServe(":8000", nil)
}
