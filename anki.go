package main

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"time"
)

type Params struct {
	Notes []Note `json:"notes"`
}

type Audio struct {
	Filename string   `json:"filename"`
	Data     string   `json:"data"`
	Fields   []string `json:"fields"`
}

type Note struct {
	DeckName  string `json:"deckName"`
	ModelName string `json:"modelName"`
	Fields    struct {
		Front      string `json:"Front"`
		Back       string `json:"Back"`
		FrontAudio string `json:"Front audio"`
	} `json:"fields"`
	Tags  []string `json:"tags"`
	Audio []Audio  `json:"audio"`
}

type Request struct {
	Action  string `json:"action"`
	Version int    `json:"version"`
	Params  struct {
		Notes []Note `json:"notes"`
	} `json:"params"`
}

// const deckName = "Default"
const deckName = "Magnus - Espanol"
const modelName = "Basic (and reversed card) with sound"

func main() {
	req := Request{
		Action:  "addNotes",
		Version: 6,
		Params: Params{
			[]Note{},
		},
	}

	// Load words from intermediate.csv
	words := []Def{}
	f, err := os.Open("intermediate.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	//r.Comma = ';'
	//i := 0
	for {
		record, err := r.Read()
		if err != nil {
			break
		}

		d := Def{
			Word:            record[0],
			FirstDefinition: record[1],
			//JoinedDefinitions: "",
			Category: Category(record[2]),
			Gender:   Gender(record[3]),
		}
		words = append(words, d)
		//i++
		//if i > 10 {
		//	break
		//}
	}

	//words := [][]string{{"mundo", "world"}}
	for _, w := range words {
		note := Note{
			DeckName:  deckName,
			ModelName: modelName,
		}
		note.Fields.Front = w.Front()
		note.Fields.Back = w.Back()
		note.Fields.FrontAudio = "" //w[0] + ".mp3"
		note.Tags = []string{
			"kindle-to-anki",
			// Current date
			time.Now().Format("2006-01-02"),
		}

		note.Audio = []Audio{
			{
				Filename: w.Word + ".mp3",
				Data:     tts(w.Word),
				Fields:   []string{"Front audio"},
			},
		}
		req.Params.Notes = append(req.Params.Notes, note)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")

	if err := enc.Encode(req); err != nil {
		panic(err)
	}
}
