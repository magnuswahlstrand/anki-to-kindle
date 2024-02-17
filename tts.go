package main

import (
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"context"
	"encoding/base64"
	"fmt"
	"log"
)

func tts() {

	ctx := context.Background()
	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Build the request for the text you want to convert to speech
	req := &texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: "Mundo"},
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
	resp, err := client.SynthesizeSpeech(ctx, req)
	if err != nil {
		log.Fatalf("Failed to synthesize speech: %v", err)
	}

	encodedStr := base64.StdEncoding.EncodeToString(resp.AudioContent)
	fmt.Println(encodedStr)

	fmt.Println("Audio content written to file: output.mp3")
}
