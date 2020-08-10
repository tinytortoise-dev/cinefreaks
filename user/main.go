package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Review Struct (Model)
type Review struct {
	ReviewID    string `json:"reviewid"`
	UserID      string `json:"userid"`
	MovieTitle  string `json:"movietitle"`
	ReviewTitle string `json:"reviewtitle"`
	Score       int    `json:"score"`
	Comment1    string `json:comment1`
	Comment2    string `json:comment2`
	Comment3    string `json:comment3`
}

// Init reviews var as a slice Review struct
var reviews []Review

func getReviews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

func getReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get params
	//Loop through reviews and find with reviewid
	for _, review := range reviews {
		if review.ReviewID == params["reviewid"] {
			json.NewEncoder(w).Encode(review)
			return
		}
	}
	json.NewEncoder(w).Encode(&Review{})
}

func createReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var review Review
	_ = json.NewDecoder(r.Body).Decode(&review)
	reviews = append(reviews, review)
	json.NewEncoder(w).Encode(review)
}

func updateReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, review := range reviews {
		if review.ReviewID == params["reviewid"] {
			reviews = append(reviews[:index], reviews[index+1:]...)
			var review Review
			_ = json.NewDecoder(r.Body).Decode(&review)
			review.ReviewID = params["reviewid"]
			reviews = append(reviews, review)
			json.NewEncoder(w).Encode(review)
			return
		}
	}
	json.NewEncoder(w).Encode(reviews)
}

func deleteReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, review := range reviews {
		if review.ReviewID == params["reviewid"] {
			reviews = append(reviews[:index], reviews[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(reviews)
}

func main() {
	// Init Router
	r := mux.NewRouter()

	// Mock Data
	reviews = append(reviews, Review{
		ReviewID:    "1",
		UserID:      "999",
		MovieTitle:  "Nuovo Cinema Paradiso",
		ReviewTitle: "Awesome!",
		Score:       99,
		Comment1:    "Toto",
		Comment2:    "Scene where Toto leaves his town and Alfredo sees him off at the station, saying 'Forget everything about this town. Live your future in the world.'",
		Comment3:    "highly recommended",
	})
	reviews = append(reviews, Review{
		ReviewID:    "2",
		UserID:      "888",
		MovieTitle:  "La Vita e Bella!",
		ReviewTitle: "Awesome!",
		Score:       98,
		Comment1:    "Guido",
		Comment2:    "Last scene. Definetly.",
		Comment3:    "highly recommended, too",
	})

	// Route Handlers / Endpoints
	r.HandleFunc("/cinefreaks/reviews", getReviews).Methods("GET")
	r.HandleFunc("/cinefreaks/reviews/{reviewid}", getReview).Methods("GET")
	r.HandleFunc("/cinefreaks/reviews", createReview).Methods("POST")
	r.HandleFunc("/cinefreaks/reviews/{reviewid}", updateReview).Methods("PUT")
	r.HandleFunc("/cinefreaks/reviews/{reviewid}", deleteReview).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))

}
