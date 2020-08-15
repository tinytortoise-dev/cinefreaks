package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Review represents a single review data
type Review struct {
	ReviewID  string
	UserID    string
	Title     string
	FilmName  string
	Comment   string
	Score     int
	IsDeleted bool
}

var reviews []Review

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/reviews/{id}", singleReviewHandler).Methods("GET")
	r.HandleFunc("/reviews/{id}", updateReviewHandler).Methods("PUT")
	r.HandleFunc("/reviews/{id}", deleteReviewHandler).Methods("DELETE")
	r.HandleFunc("/reviews", reviewHandler).Methods("GET")
	r.HandleFunc("/reviews", createReviewHandler).Methods("POST")
	http.Handle("/", r)

	fmt.Println("updated to latest")
	fmt.Println("review service server started on port 8001")
	log.Fatal(http.ListenAndServe(":8001", nil))
}

func singleReviewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	givenId := vars["id"]
	for _, review := range reviews {
		if review.ReviewID == givenId {
			res, err := json.Marshal(review)
			if err != nil {
				w.Write([]byte("error when marshaling"))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(res)
			return
		}
	}
	w.Write([]byte("no such review found"))
	return
}

func updateReviewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	givenId := vars["id"]
	var target Review
	for i, review := range reviews {
		if review.ReviewID == givenId {
			err := json.NewDecoder(r.Body).Decode(&target)
			if err != nil {
				w.Write([]byte("error when unmarshaling"))
				return
			}
			reviews[i].ReviewID = givenId
			reviews[i].UserID = target.UserID
			reviews[i].Title = target.Title
			reviews[i].FilmName = target.FilmName
			reviews[i].Comment = target.Comment
			reviews[i].Score = target.Score
			reviews[i].IsDeleted = target.IsDeleted
			w.Write([]byte("review updated"))
			return
		}
	}
	w.Write([]byte("review not found"))
	return
}

func deleteReviewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	givenId := vars["id"]
	for i, review := range reviews {
		if review.ReviewID == givenId {
			reviews[i].IsDeleted = true
			w.Write([]byte("review is logically deleted"))
			return
		}
	}
	w.Write([]byte("review not found"))
	return
}

func reviewHandler(w http.ResponseWriter, r *http.Request) {
	res, err := json.Marshal(reviews)
	if err != nil {
		w.Write([]byte("error when marshaling"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
	return
}

func createReviewHandler(w http.ResponseWriter, r *http.Request) {
	uuid := uuid.New()
	reviewId := uuid.String()

	new := Review{}
	err := json.NewDecoder(r.Body).Decode(&new)
	if err != nil {
		w.Write([]byte("error when parsing body"))
		return
	}
	new.ReviewID = reviewId
	reviews = append(reviews, new)
	// pass reviewId to user service
	url := fmt.Sprintf("http://user-clusterip-srv:8000/users/%s/reviewid/%s", new.UserID, reviewId)
	fmt.Println(url)
	resp, err := http.Post(url, "application/json", nil)
	defer resp.Body.Close()
	if err != nil {
		w.Write([]byte("error when sending reviewId to user service"))
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.Write([]byte("error when reading response body from user service"))
		return
	}
	fmt.Println(string(body))

	// if reviewId was successfully added to a user, ok response
	w.Write([]byte("review created"))
	return
}
