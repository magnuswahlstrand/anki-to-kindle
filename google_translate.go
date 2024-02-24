package main

import (
	"cloud.google.com/go/translate"
	"context"
	"golang.org/x/text/language"
	"log"
)

func GoogleTranslate(stemmedWords []string) [][]string {
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

	words := [][]string{}
	for i := 0; i < len(stemmedWords); i += batchSize {
		// Translate text
		translations, err := client.Translate(ctx, stemmedWords[i:min(i+batchSize, len(stemmedWords))], targetLang, nil)
		if err != nil {
			log.Fatalf("Failed to translate words: %v", err)
		}

		// Print the translations
		for j, t := range translations {
			t := t
			words = append(words, []string{stemmedWords[i+j], t.Text})
		}
	}
	return words
}
