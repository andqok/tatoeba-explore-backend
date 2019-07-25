package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

func readSentences(connStr string) {
	//connStr := "user=postgres password=ajoutee dbname=tatoeba_explore sslmode=disable"
	csvFile, err := os.Open("./data/sentences.csv")
	if err != nil {
		log.Println(err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	reader.Comma = '\t' // Use tab-delimited instead of comma <---- here!
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	csvData, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sqlStatement := `
    INSERT INTO sentences (sentence_number, lang, sentence_text)
		values ( $1,  $2,  $3)
`
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	for _, each := range csvData {
		fmt.Println(each[0], each[1], each[2])
		if len(each[1]) == 3 {
			_, err = db.Exec(sqlStatement, each[0], each[1], each[2])
			if err != nil {
				panic(err)
			}
		}
	}
}

func importLinks() {
	connStr := "user=postgres password=ajoutee dbname=tatoeba_explore sslmode=disable"
	csvFile, err := os.Open("./data/links.csv")
	if err != nil {
		log.Println(err)
	}
	defer csvFile.Close()
	reader := csv.NewReader(csvFile)

	reader.Comma = '\t' // Use tab-delimited instead of comma <---- here!
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	csvData, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	insertStatement := `INSERT INTO links (link_1, link_2) values ( $1, $2 )`
	for _, line := range csvData {
		// fmt.Println(line[0], line[1])
		_, err = db.Exec(insertStatement, line[0], line[1])
		if err != nil {
			panic(err)
		}
	}
}

func importLinks2() {
	connStr := "user=postgres password=ajoutee dbname=tatoeba_explore sslmode=disable"
	csvFile, err := os.Open("./data/links2.csv")
	if err != nil {
		log.Println(err)
	}
	defer csvFile.Close()
	reader := csv.NewReader(csvFile)

	reader.Comma = '\t' // Use tab-delimited instead of comma <---- here!
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	csvData, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	txn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := txn.Prepare(pq.CopyIn("links", "link_1", "link_2"))
	if err != nil {
		log.Fatal(err)
	}
	// insertStatement := `INSERT INTO links (link_1, link_2) values ( $1, $2 )`
	for _, line := range csvData {
		// fmt.Println(line[0], line[1])
		// _, err = db.Exec(insertStatement, line[0], line[1])
		// if err != nil {
		// 	panic(err)
		// }
		_, err = stmt.Exec(line[0], line[1])
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	// importLinks()
	// runServer()
	importLinks2()
}
