package main

import (
	"database/sql"
	"encoding/csv"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

const vocabDBPath = "db/vocab.db"
const vocabLanguage = "es" // Set the vocabLanguage you're interested in

type Translation struct {
	Word        string `json:"word"`
	Translation string `json:"translation"`
}

func main() {
	//tts()
	//return
	// Open the SQLite database
	db, err := sql.Open("sqlite3", vocabDBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Prepare SQL query to select stem words for the specified vocabLanguage
	query := `SELECT DISTINCT stem FROM WORDS WHERE lang = ? LIMIT ?`
	rows, err := db.Query(query, vocabLanguage, limit)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the rows and print each stem word
	stemmedWords := []string{}
	for rows.Next() {
		var stem string
		if err := rows.Scan(&stem); err != nil {
			log.Fatal(err)
		}
		stemmedWords = append(stemmedWords, stem)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	//words := GoogleTranslate(stemmedWords)
	words := Translate2(stemmedWords)

	f, err := os.OpenFile("intermediate.csv", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	w := csv.NewWriter(f)
	if err := w.WriteAll(words); err != nil {
		log.Fatal(err)
	}
}

const limit = 1000
const batchSize = 128
