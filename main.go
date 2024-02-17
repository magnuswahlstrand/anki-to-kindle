package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/text/language"
	"log"

	"cloud.google.com/go/translate"
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
	query := `SELECT DISTINCT stem FROM WORDS WHERE lang = ? LIMIT 10`
	rows, err := db.Query(query, vocabLanguage)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the rows and print each stem word
	fmt.Println("Stem Words:")
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

	ctx := context.Background()

	// Initialize Google Cloud Translation client
	client, err := translate.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Set the target vocabLanguage
	targetLang, err := language.Parse("en")
	if err != nil {
		log.Fatalf("Failed to parse target vocabLanguage: %v", err)
	}

	words := []Translation{}
	for i := 0; i < len(stemmedWords); i += 128 {
		// Translate text
		translations, err := client.Translate(ctx, stemmedWords[i:min(i+128, len(stemmedWords))], targetLang, nil)
		if err != nil {
			log.Fatalf("Failed to translate words: %v", err)
		}

		// Print the translations
		for i, translation := range translations {
			words = append(words, Translation{Word: stemmedWords[i], Translation: translation.Text})
		}
	}
}
