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
	words := [][]string{}
	f, err := os.Open("intermediate.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	//i := 0
	for {
		record, err := r.Read()
		if err != nil {
			break
		}
		words = append(words, record)
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
		note.Fields.Front = w[0]
		note.Fields.Back = w[1]
		note.Fields.FrontAudio = "" //w[0] + ".mp3"
		note.Tags = []string{
			"kindle-to-anki",
			// Current date
			time.Now().Format("2006-01-02"),
		}

		note.Audio = []Audio{
			{
				Filename: w[0] + ".mp3",
				Data:     tts(w[0]),
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
