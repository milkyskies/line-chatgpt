package line

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"cloud.google.com/go/translate"
	"github.com/hajimehoshi/go-mp3"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var (
	ErrSendMessageFailed          = errors.New("failed to send message")
	ErrReceiveMessageNotSupported = errors.New("receive message not supported for LINE")
)

func (l *Chat) SendMessage(userID string, message string) error {
	_, err := l.Client.PushMessage(userID, linebot.NewTextMessage(message)).Do()
	if err != nil {
		return ErrSendMessageFailed
	}
	return nil
}

func (l *Chat) SendAudioMessage(userID string, id string) error {
	f, err := os.Open(filepath.Join("content/whisper/audio_replies", fmt.Sprintf("%s.mp3", id)))
	if err != nil {
		return err
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		return err
	}

	sampleSize := 4
	samples := int(d.Length()) / sampleSize
	duration := samples / d.SampleRate() * 1000

	fmt.Println("Duration: ", duration)

	hostname := os.Getenv("HOSTNAME")

	audioMessage := linebot.NewAudioMessage(hostname+"/audio_replies/"+id, duration)

	if _, err := l.Client.PushMessage(userID, audioMessage).Do(); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// TODO: Move this
func (l *Chat) GenerateAudio(message string, id string) error {
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
