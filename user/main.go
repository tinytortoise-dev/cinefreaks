package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tinytortoise-dev/cinefreaks/user/helper"
)

// User represents a single user information
type User struct {
	UserID      string   `json:"userId"`
	MailAddress string   `json:"mailAddress"`
	Password    string   `json:"password"`
	ReviewIds   []string `json:"reviewIds"` // I can use ref/populate mongo feature?
}

var users []User

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/users/{id}/reviewid/{reviewId}", addReviewToUser).Methods("POST") // need verification
	r.HandleFunc("/users/{id}/reviews", getUserReviews).Methods("GET")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	http.Handle("/", r)

	fmt.Println("user service server started on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))

}

func addReviewToUser(w http.ResponseWriter, r *http.Request) {
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
		helper.NotFound(w)
		return
	}
	w.Write([]byte("review added"))
}

func getUserReviews(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]
	// currently, return review ids. in the future, return the slices of actual review struct
	for _, u := range users {
		if userId == u.UserID {
			res, err := json.Marshal(u.ReviewIds)
			if err != nil {
				helper.ServerError(w)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(res)
			return
		}
	}
	helper.NotFound(w)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println(id)
	var user User
	foundUser := false
	for _, u := range users {
		if id == u.UserID {
			user = u
			foundUser = true
		}
	}
	if !foundUser {
		helper.NotFound(w)
		return
	}
	res, err := json.Marshal(user)
	if err != nil {
		helper.ServerError(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	res, err := json.Marshal(users)
	if err != nil {
		helper.ServerError(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/users" {
	// 	helper.NotFound(w)
	// 	return
	// }

	// if r.Method == "GET" {
	// 	res, err := json.Marshal(users)
	// 	if err != nil {
	// 		helper.ServerError(w)
	// 		return
	// 	}
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.Write(res)
	// }

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		helper.ServerError(w)
		return
	}
	data, err := duplicateUser(users, user)
	if err != nil {
		res := helper.ErrorJsons{}
		res.AddMessageAndDataEntry(err.Error(), data)
		res.JsonError(w, http.StatusBadRequest)
		return
	}
	users = append(users, user)
	w.Write([]byte("user created"))
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
