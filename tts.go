package main

import (
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"context"
	"encoding/base64"
	"log"
)

var client *texttospeech.Client

func init() {
	c, err := texttospeech.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	client = c
}

func tts(word string) string {
	// Build the request for the text you want to convert to speech
	req := &texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: word},
		},
		// Specify the voice and language code
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "es-ES",
			Name:         "es-ES-Standard-A", // Adjust the voice name as needed
		},
		// Specify the audio format you want to receive
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},
	}

	// Perform the text-to-speech request on the text input with the selected voice parameters and audio file type
	resp, err := client.SynthesizeSpeech(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to synthesize speech: %v", err)
	}

	return base64.StdEncoding.EncodeToString(resp.AudioContent)
}
