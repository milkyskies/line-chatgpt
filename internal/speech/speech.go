package speech

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"cloud.google.com/go/translate"
)

func GenerateAudio(message string, id string) error {
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	language, err := detectLanguage(message)
	if err != nil {
		return err
	}

	langTag := language.Language.String()

	var voiceSelectionParams texttospeechpb.VoiceSelectionParams
	switch langTag {
	case "ja":
		voiceSelectionParams = texttospeechpb.VoiceSelectionParams{
			LanguageCode: "ja-JP",
			Name:         "ja-JP-Neural2-C",
		}
	case "en":
		voiceSelectionParams = texttospeechpb.VoiceSelectionParams{
			LanguageCode: "en-US",
			Name:         "en-US-Standard-A",
		}
	default:
		voiceSelectionParams = texttospeechpb.VoiceSelectionParams{
			LanguageCode: langTag,
			SsmlGender:   texttospeechpb.SsmlVoiceGender_MALE,
		}
	}

	req := texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: message},
		},
		Voice: &voiceSelectionParams,
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		return fmt.Errorf("could not synthesize speech: %f", err)
	}

	outputDir := "content/whisper/audio_replies"

	if err = os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	outputFilePathMP3 := filepath.Join(outputDir, id+".mp3")
	outputFilePathM4A := filepath.Join(outputDir, id+".m4a")

	if err := ioutil.WriteFile(outputFilePathMP3, resp.AudioContent, 0644); err != nil {
		return err
	}
	fmt.Printf("Audio content written to file: %v\n", outputFilePathMP3)

	err = convertMP3ToM4A(outputFilePathMP3, outputFilePathM4A)

	return err
}

func convertMP3ToM4A(inputFile, outputFile string) error {
	cmd := exec.Command("ffmpeg", "-i", inputFile, "-c:a", "aac", outputFile)
	return cmd.Run()
}

func detectLanguage(text string) (*translate.Detection, error) {
	ctx := context.Background()
	client, err := translate.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("translate.NewClient: %v", err)
	}
	defer client.Close()
	lang, err := client.DetectLanguage(ctx, []string{text})
	if err != nil {
		return nil, fmt.Errorf("DetectLanguage: %v", err)
	}
	if len(lang) == 0 || len(lang[0]) == 0 {
		return nil, fmt.Errorf("DetectLanguage return value empty")
	}
	return &lang[0][0], nil
}
