package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var connStr = "user=postgres password=ajoutee dbname=tatoeba_explore sslmode=disable"

// Sentence is my comment
type Sentence struct {
	ID     int    `json:"id"`
	Number int    `json:"number"`
	Text   string `json:"text"`
	Lang   string `json:"lang"`
}

func getServer(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	number := params["sentence_number"]
	// fmt.Println(number)
	// fmt.Fprint(w, "Hello")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	row := db.QueryRow("SELECT * FROM sentences WHERE sentence_number = $1", number)
	// fmt.Println(row)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	snt := Sentence{}
	// fmt.Println(snt)
	err2 := row.Scan(
		&snt.ID,
		&snt.Number,
		&snt.Text,
		&snt.Lang)
	if err2 != nil {
		// fmt.Println(err2)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	fmt.Println(snt)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snt)
}

func runServer() {
	router := mux.NewRouter()
	router.HandleFunc("/sentence/{sentence_number}", getServer).Methods("GET")
	log.Fatal(http.ListenAndServe(":3001", router))
}
