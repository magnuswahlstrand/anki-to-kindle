package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// https://www.dictionaryapi.com/products/json
type RawDefinition struct {
	Meta struct {
		ID        string   `json:"id"`
		UUID      string   `json:"uuid"`
		Lang      string   `json:"lang"`
		Sort      string   `json:"sort"`
		Src       string   `json:"src"`
		Section   string   `json:"section"`
		Stems     []string `json:"stems"`
		Offensive bool     `json:"offensive"`
	} `json:"meta"`
	Hom int `json:"hom"`
	Hwi struct {
		Hw  string `json:"hw"`
		Prs []struct {
			Sound struct {
				Audio string `json:"audio"`
			} `json:"sound"`
		} `json:"prs"`
	} `json:"hwi"`
	Fl  string `json:"fl"`
	Def []struct {
		Sseq [][][]any `json:"sseq"`
		Vd   string    `json:"vd,omitempty"`
	} `json:"def"`
	Suppl struct {
		Cjts []struct {
			Cjid string   `json:"cjid"`
			Cjfs []string `json:"cjfs"`
		} `json:"cjts"`
	} `json:"suppl"`
	Shortdef []string `json:"shortdef"`
}

const dictionaryUrl = "https://www.dictionaryapi.com/api/v3/references/spanish/json/"

type Category string

const (
	Adjective    Category = "adjective"
	Adverb                = "adverb"
	Conjunction           = "conjunction"
	Interjection          = "interjection"
	Noun                  = "noun"
	Preposition           = "preposition"
	Pronoun               = "pronoun"
	Verb                  = "verb"
	Unknown               = "Unknown"
)

type Gender string

const (
	Feminine  Gender = "feminine"
	Masculine        = "masculine"
	NA               = "N/A"
)

type Def struct {
	Word              string
	FirstDefinition   string
	JoinedDefinitions string
	Category          Category
	Gender            Gender
}

var baseNounColor = "200, 200, 200"

func (d Def) Front() string {
	var prefix string
	var nounColor string
	if d.Category == Noun {
		switch d.Gender {
		case "feminine":
			prefix = "la "
			nounColor = "255, 105, 180"
		case "masculine":
			prefix = "el "
			nounColor = "0, 178, 255"
		}
	}

	return colorize(prefix+d.Word, d.Category, nounColor)
}

func colorize(word string, category Category, nounColor string) string {
	switch category {
	case Adjective:
		return addColor(word, "255, 225, 25")
	case Adverb:
		return addColor(word, "0, 217, 217")
	case Verb:
		return addColor(word, "255, 127, 80")
	case Noun:
		return addColor(word, nounColor)
	default:
		return word
	}
}

func (d Def) Back() string {
	return colorize(d.FirstDefinition, d.Category, baseNounColor)
}

func addColor(s, color string) string {
	return "<span style=\"color: rgb(" + color + ")\">" + s + "</font>"
}

func getWordDefinition(word string) (Def, error) {
	fmt.Println(word)
	resp, err := http.Get(dictionaryUrl + word + "?key=" + os.Getenv("DICTIONARY_API_KEY"))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var definitions []RawDefinition
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Def{}, err
	}

	if err := json.Unmarshal(bodyBytes, &definitions); err != nil {
		return Def{}, fmt.Errorf("Failed to unmarshal response: %v, %s", err, string(bodyBytes))
	}

	if len(definitions) == 0 {
		return Def{}, fmt.Errorf("No definition found for %s", word)
	}

	topDefinition := definitions[0]

	if len(topDefinition.Shortdef) == 0 {
		return Def{}, fmt.Errorf("No short definition found for %s", word)
	}

	d := Def{
		Word:              topDefinition.Hwi.Hw,
		Category:          parseCategory(topDefinition.Fl),
		Gender:            parseGender(topDefinition.Fl),
		FirstDefinition:   topDefinition.Shortdef[0],
		JoinedDefinitions: strings.Join(topDefinition.Shortdef, ", "),
	}

	return d, nil
}

func parseGender(fl string) Gender {
	switch {
	case strings.Contains(fl, "feminine"):
		return Feminine
	case strings.Contains(fl, "masculine"):
		return Masculine
	default:
		return NA
	}
}

func parseCategory(fl string) Category {
	switch {
	case strings.Contains(fl, "adjective"):
		return Adjective
	case strings.Contains(fl, "adverb"):
		return Adverb
	case strings.Contains(fl, "conjunction"):
		return Conjunction
	case strings.Contains(fl, "interjection"):
		return Interjection
	case strings.Contains(fl, "noun"):
		return Noun
	case strings.Contains(fl, "preposition"):
		return Preposition
	case strings.Contains(fl, "pronoun"):
		return Pronoun
	case strings.Contains(fl, "verb"):
		return Verb
	default:
		return Unknown
	}
}

func Translate2(stemmedWords []string) [][]string {
	words := [][]string{}
	for _, word := range stemmedWords {
		def, err := getWordDefinition(word)
		if err != nil {
			fmt.Println(err)
			continue
		}
		words = append(words, []string{def.Word, def.FirstDefinition, string(def.Category), string(def.Gender), def.JoinedDefinitions})
	}
	return words
}
