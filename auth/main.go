package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/twinj/uuid"
)

var client *redis.Client

func init() {
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	client = redis.NewClient(&redis.Options{
		Addr: dsn,
	})

	_, err := client.Ping().Result()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

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

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
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
	td, err := createToken(incoming.UserId)
	if err != nil {
		// error handling
		w.Write([]byte("61"))
		return
	}

	saveErr := CreateAuth(incoming.UserId, td)
	if saveErr != nil {
		// error handling
		w.Write([]byte("61"))
		return
	}

	// send the token back to client
	var m map[string]string
	m = make(map[string]string)
	m["access_token"] = td.AccessToken
	m["refresh_token"] = td.RefreshToken
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

func CreateAuth(userId string, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) // unix to UTC
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := client.Set(td.AccessUuid, userId, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}

	errRefresh := client.Set(td.RefreshUuid, userId, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}

	return nil
}

func createToken(userId string) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userId
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACEESS_SECRET")))
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userId
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	return td, nil

	// atClaims := jwt.MapClaims{}
	// atClaims["authorized"] = true
	// atClaims["user_id"] = userId
	// atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	// at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	// token, err := at.SignedString([]byte(os.Getenv("JWT_KEY")))
	// if err != nil {
	// 	// error handling
	// 	return "", err
	// }
	// return token, nil
}
