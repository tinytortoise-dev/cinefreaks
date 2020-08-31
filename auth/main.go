package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/auth/signin", signin).Methods("POST")
	http.Handle("/", r)

	fmt.Println("auth service server started on port 8002")
	log.Fatal(http.ListenAndServe(":8002", nil))
}

type UserCredentials struct {
	UserId   string
	Password string
}

func signin(w http.ResponseWriter, r *http.Request) {
	var incoming UserCredentials
	err := json.NewDecoder(r.Body).Decode(&incoming)
	if err != nil {
		// error handling
		w.Write([]byte("failed to decode body"))
		return
	}
	// get user info by that userId from user service
	url := fmt.Sprintf("http://user-clusterip-srv:8000/users/%s", incoming.UserId)
	resp, err := http.Get(url)
	if err != nil {
		// error handling
		w.Write([]byte("39"))
		return
	}
	// check if userId and password are correct
	// if not, return
	var database UserCredentials
	err = json.NewDecoder(resp.Body).Decode(&database)
	defer resp.Body.Close()
	if err != nil {
		// error handling
		w.Write([]byte(err.Error()))
		return
	}
	if database.UserId != incoming.UserId ||
		database.Password != incoming.Password {
		// error handling
		w.Write([]byte("54"))
		return
	}
	// create token by userId
	token, err := createToken(incoming.UserId)
	if err != nil {
		// error handling
		w.Write([]byte("61"))
		return
	}
	// send the token back to client
	var m map[string]string
	m = make(map[string]string)
	m["token"] = token
	json, err := json.Marshal(m)
	if err != nil {
		// error handling
		w.Write([]byte("70"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
	return
}

func createToken(userId string) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		// error handling
		return "", err
	}
	return token, nil
}
