package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const connStr = "user=postgres password=ajoutee dbname=tatoeba_explore sslmode=disable"
const delim = "?!»«()/:.;-,*—"

func removePunctuation(s string) string {
	return strings.Map(
		func(r rune) rune {
			if strings.Contains(delim, string(r)) {
				return -1
			}
			return r
		},
		s,
	)
}

// Sentence my-f-ckomment
type Sentence struct {
	ID     int    `json:"id"`
	Number int    `json:"number"`
	Text   string `json:"text"`
	Lang   string `json:"lang"`
}

// WordFreq a word frequency object
type WordFreq struct {
	Word      string `json:"word"`
	Frequency int    `json:"frequency"`
}

/* Get specific sentence by its number */
func getSentence(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	number := params["sentence_number"]
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	row := db.QueryRow("SELECT * FROM sentences WHERE sentence_number = $1", number)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	snt := Sentence{}
	err2 := row.Scan(
		&snt.ID,
		&snt.Number,
		&snt.Text,
		&snt.Lang)
	if err2 != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	fmt.Println(snt)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snt)
}

func getSplittedWords() {
	currLang := "eng"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	rows, err1 := db.Query("SELECT * FROM sentences WHERE lang = $1", currLang)
	if err1 != nil {
		fmt.Println("Error happened", err1)
		return
	}
	defer rows.Close()

	insertStatement := `INSERT INTO words (word, sentence_number, lang) VALUES ($1, $2, $3)`
	for rows.Next() {
		snt := Sentence{}
		err2 := rows.Scan(
			&snt.ID,
			&snt.Number,
			&snt.Text,
			&snt.Lang)
		if err2 != nil {
			fmt.Println("Error happened")
			return
		}
		for _, word := range strings.Fields(snt.Text) {
			lower := strings.ToLower(word)
			removed := removePunctuation(lower)
			fmt.Println(lower, removed)
			_, err = db.Exec(insertStatement, removed, snt.Number, currLang)
			if err != nil {
				panic(err)
			}
		}
	}
}

func getTopWords(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	limit := params["limit"]
	lang := params["lang"]
	queryStr := `SELECT word, COUNT(*) AS frequency
		FROM words
		WHERE lang = $1
		GROUP BY word
		ORDER BY COUNT(*) DESC
		LIMIT $2;`

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	rows, err1 := db.Query(queryStr, lang, limit)
	if err1 != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	resp := []WordFreq{}
	for rows.Next() {
		wordFreq := WordFreq{}
		err2 := rows.Scan(
			&wordFreq.Word,
			&wordFreq.Frequency,
		)
		if err2 != nil {
			fmt.Println("Error happened")
			return
		}
		resp = append(resp, wordFreq)
	}
	fmt.Println(rows)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func searchText(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	text := params["text"]
	lang := params["lang"]
	queryStr := `select * from sentences where sentence_text ilike '%` + text + `%' and lang = $1;`
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	rows, err1 := db.Query(queryStr, lang)
	if err1 != nil {
		fmt.Println(err1)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	fmt.Println(rows)
	sentences := []Sentence{}
	for rows.Next() {
		snt := Sentence{}
		err2 := rows.Scan(
			&snt.ID,
			&snt.Number,
			&snt.Text,
			&snt.Lang)
		if err2 != nil {
			fmt.Println("Error happened")
			return
		}
		sentences = append(sentences, snt)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sentences)
}

func runServer() {
	router := mux.NewRouter()
	router.HandleFunc("/sentence/{sentence_number}", getSentence).Methods("GET")
	router.HandleFunc("/word/top/{limit}/{lang}", getTopWords).Methods("GET")
	router.HandleFunc("/search/{lang}/{text}", searchText).Methods("GET")
	log.Fatal(http.ListenAndServe(":3001", router))
}
