package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// User represents a single user information
type User struct {
	UserID      string   `json:"userId"`
	MailAddress string   `json:"mailAddress"`
	Password    string   `json:"password"`
	ReviewIds   []string `json:"reviewIds"`
}

var users []User

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/users/{id}/reviewid/{reviewId}", addUserReviewHandler)
	r.HandleFunc("/users/{id}/reviews", userReviewsHandler)
	r.HandleFunc("/users/{id}", singleUserHandler)
	r.HandleFunc("/users", userHandler)
	http.Handle("/", r)

	fmt.Println("user service server started on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))

}

func addUserReviewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("addUserReviewHandler was called")
	vars := mux.Vars(r)
	userId := vars["id"]
	reviewId := vars["reviewId"]
	var user User
	foundUser := false
	for i, u := range users {
		if userId == u.UserID {
			user = u
			foundUser = true
			user.ReviewIds = append(user.ReviewIds, reviewId)
			users[i] = user
		}

	}
	if !foundUser {
		w.Write([]byte("no such user found"))
		return
	}
	w.Write([]byte("review added"))
}

func userReviewsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]
	// currently, return review ids. in the future, return the slices of actual review struct
	for _, u := range users {
		if userId == u.UserID {
			// ids := struct{ reviewIds []int }{u.ReviewIds}
			res, err := json.Marshal(u.ReviewIds)
			if err != nil {
				w.Write([]byte("error when marshaling"))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(res)
			return
		}
	}
	w.Write([]byte("no such user found"))
}

func singleUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var user User
	foundUser := false
	for _, u := range users {
		if id == u.UserID {
			user = u
			foundUser = true
		}
	}
	if !foundUser {
		w.Write([]byte("no such user found"))
		return
	}
	res, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "error when marshaling a single user struct", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/users" {
		http.Error(w, "404 Not found", http.StatusNotFound)
		return
	}

	if r.Method == "GET" {
		res, err := json.Marshal(users)
		if err != nil {
			http.Error(w, "error when marshaling user structs", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	}

	if r.Method == "POST" {
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "error when reaing body", http.StatusBadRequest)
			return
		}
		data, err := duplicateUser(users, user)
		if err != nil {
			res := ErrorJsons{}
			res.addMessageAndDataEntry(err.Error(), data)
			res.jsonError(w, http.StatusBadRequest)
			return
		}
		users = append(users, user)
		fmt.Println("created a new user")
		fmt.Println(user)
		fmt.Fprintf(w, "user was created")
	}

}

func duplicateUser(users []User, user User) (string, error) {
	for _, u := range users {
		if u.UserID == user.UserID {
			return user.UserID, errors.New("This userId is already in use")
		}
		if u.MailAddress == user.MailAddress {
			return user.MailAddress, errors.New("This mailAddress is already in use")
		}
	}
	return "", nil
}

func (e *ErrorJsons) addMessageAndDataEntry(message, data string) {
	errJson := ErrorJson{}
	errJson.setMessage(message)
	errJson.setData(data)
	e.addErrorJson(errJson)
}

func (e *ErrorJsons) addMessageEntry(message string) {
	errJson := ErrorJson{}
	errJson.setMessage(message)
	e.addErrorJson(errJson)
}

func (e *ErrorJsons) jsonError(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	b, err := json.Marshal(e)
	if err != nil {
		http.Error(w, "error when marshaling struct", http.StatusInternalServerError)
	}
	// json.NewEncoder(w).Encode(err)
	http.Error(w, string(b), code)
}

func (e *ErrorJson) setMessage(message string) {
	e.Message = message
}

func (e *ErrorJson) setData(data string) {
	e.Data = data
}

// func (e ErrorJson) getErrorJson() ErrorJson {
// 	return e
// }

func (e *ErrorJsons) addErrorJson(errorJson ErrorJson) {
	e.Errors = append(e.Errors, errorJson)
}

// func (e ErrorJsons) getErrorJsons() ErrorJsons {
// 	return e
// }

type ErrorJson struct {
	Message string `json:"message"` // no space between json: and value
	Data    string `json:"data"`    // optional
}

type ErrorJsons struct {
	Errors []ErrorJson `json:"errors"`
}
