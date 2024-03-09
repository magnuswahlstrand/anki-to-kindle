# Kindle to Anki
Transfers words from the Kindle to Anki. It is configured to work with a Spanish dictionary, but can be adapted to any language. 
Currently it takes 3 steps to upload to Anki

## Prerequisites
You need to have a dictionary API key from [Merriam-Webster](https://dictionaryapi.com/), and set the environment variable `DICTIONARY_API_KEY` to your key.
```
DICTIONARY_API_KEY=your_key
```
You will also need GCP credentials for the TTS API.

## Extract words from Kindle
1. Copy vocab.db from Kindle to your computer
2. `go run ./main.go ./dictionary.go`

This creates a file called intermediate.csv, with words and definitions.

## Create upload format
1. go run ./anki.go ./tts.go ./dictionary.go | pbcopy
2. Paste into upload.http

## Upload to Anki
1. Open Anki
2. Make sure [AnkiConnect](https://ankiweb.net/shared/info/2055492159) is installed. This exposes a REST API.
3. Open upload.http and send request #2


# TODO
* [x] Look up if word is a verb or noun
* [x] Add colors to verbs/nouns
* [ ] Handle darkmode in anki
* [ ] Define global classes in Anki?
* [ ] Mark verbs as -er, -ir, -ar
* [ ] Source words from a file
