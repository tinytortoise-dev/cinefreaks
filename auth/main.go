package main

import (
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

func signin(w http.ResponseWriter, r *http.Request) {
	// extract userId and password
	// get user info by that userId from user service
	// check if userId and password are correct
	// if not, return
	// create token by userId
	// send the token back to client
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
