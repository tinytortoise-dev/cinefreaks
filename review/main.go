package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tinytortoise-dev/cinefreaks/review/helper"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Review represents a single review data
type Review struct {
	ReviewID  string `json:"reviewId"`
	UserID    string `json:"userId"`
	Title     string `json:"title"`
	FilmName  string `json:"filmName"`
	Comment   string `json:"comment"`
	Score     int    `json:"score"`
	IsDeleted bool   `json:"isDeleted"`
}

var reviews []Review

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/reviews/{id}", getReviewByReviewId).Methods("GET")
	r.HandleFunc("/reviews/{id}", updateReviewByReviewId).Methods("PUT")
	r.HandleFunc("/reviews/{id}", deleteReviewByReviewId).Methods("DELETE")
	r.HandleFunc("/reviews", getReviews).Methods("GET")
	r.HandleFunc("/reviews", addReview).Methods("POST")
	http.Handle("/", r)

	fmt.Println("review service server started on port 8001")
	log.Fatal(http.ListenAndServe(":8001", nil))
}

func getReviewByReviewId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	givenId := vars["id"]
	for _, review := range reviews {
		if review.ReviewID == givenId {
			res, err := json.Marshal(review)
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
	return
}

func updateReviewByReviewId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	givenId := vars["id"]
	var target Review
	for i, review := range reviews {
		if review.ReviewID == givenId {
			err := json.NewDecoder(r.Body).Decode(&target)
			if err != nil {
				helper.ServerError(w)
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
	helper.NotFound(w)
	return
}

func deleteReviewByReviewId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	givenId := vars["id"]
	for i, review := range reviews {
		if review.ReviewID == givenId {
			reviews[i].IsDeleted = true
			w.Write([]byte("review is logically deleted"))
			return
		}
	}
	helper.NotFound(w)
	return
}

func getReviews(w http.ResponseWriter, r *http.Request) {
	res, err := json.Marshal(reviews)
	if err != nil {
		helper.ServerError(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
	return
}

func addReview(w http.ResponseWriter, r *http.Request) {
	uuid := uuid.New()
	reviewId := uuid.String()

	new := Review{}
	err := json.NewDecoder(r.Body).Decode(&new)
	if err != nil {
		helper.ServerError(w)
		return
	}
	new.ReviewID = reviewId
	reviews = append(reviews, new)
	// pass reviewId to user service
	url := fmt.Sprintf("http://user-clusterip-srv:8000/users/%s/reviewid/%s", new.UserID, reviewId)
	resp, err := http.Post(url, "application/json", nil)
	defer resp.Body.Close()
	if err != nil {
		helper.ServerError(w)
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
