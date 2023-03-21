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
    ErrSendMessageFailed = errors.New("failed to send message")
    ErrReceiveMessageNotSupported = errors.New("receive message not supported for LINE")
)

func (l *LineChat) SendMessage(userID string, message string) error {
    _, err := l.Client.PushMessage(userID, linebot.NewTextMessage(message)).Do()
    if err != nil {
        return ErrSendMessageFailed
    }
    return nil
}

func (l *LineChat) SendAudioMessage(userID string, id string) error {
    f, err := os.Open(filepath.Join("content/whisper/audio_replies", fmt.Sprintf("%s.mp3", id)),)
    if err != nil {
        return err
    }
    defer f.Close()

    // Decode the MP3 file
    d, err := mp3.NewDecoder(f)
    if err != nil {
        return err
    }

    sampleSize := 4                    // From documentation.
    samples := int(d.Length()) / sampleSize      // Number of samples.
    duration := samples / d.SampleRate() * 1000 // Audio length in seconds.

    fmt.Println("Duration: ", duration)

    audioMessage := linebot.NewAudioMessage("https://api.yozora.dev/audio_replies/" + id, duration)
    
    if _, err := l.Client.PushMessage(userID, audioMessage).Do(); err != nil {
        fmt.Println(err)
        return err
    }

    return nil
}

func (l *LineChat) ReceiveMessage(userID string) (string, error) {
    // Pass the message somewhere else to handle

    return "", ErrReceiveMessageNotSupported
}

// TODO: Move this
func (l *LineChat) GenerateAudio(message string, id string) error {
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
            Name: "ja-JP-Neural2-C",
        }
    case "en":
        voiceSelectionParams = texttospeechpb.VoiceSelectionParams{
            LanguageCode: "en-US",
            Name: "en-US-Standard-A",
        }
    default:
        voiceSelectionParams = texttospeechpb.VoiceSelectionParams{
            LanguageCode: langTag,
            SsmlGender: texttospeechpb.SsmlVoiceGender_MALE,
        }
    }

    // Perform the text-to-speech request on the text input with the selected
    // voice parameters and audio file type.
    req := texttospeechpb.SynthesizeSpeechRequest{
            // Set the text input to be synthesized.
            Input: &texttospeechpb.SynthesisInput{
                    InputSource: &texttospeechpb.SynthesisInput_Text{Text: message},
            },
            // Build the voice request, select the language code ("en-US") and the SSML
            // voice gender ("neutral").
            Voice: &voiceSelectionParams,
            // Select the type of audio file you want returned.
            AudioConfig: &texttospeechpb.AudioConfig{
                    AudioEncoding: texttospeechpb.AudioEncoding_MP3,
            },
    }

    resp, err := client.SynthesizeSpeech(ctx, &req)
    if err != nil {
            log.Fatal(err)
    }

    // The resp's AudioContent is binary
    
    outputDir := "content/whisper/audio_replies"

	if err = os.MkdirAll(outputDir, 0755); err != nil {
        return err
    }

	outputFilePathMP3 := filepath.Join(outputDir, id + ".mp3")
    outputFilePathM4A := filepath.Join(outputDir, id+".m4a")

	if err := ioutil.WriteFile(outputFilePathMP3, resp.AudioContent, 0644); err != nil {
        return err
    }
	fmt.Printf("Audio content written to file: %v\n", outputFilePathMP3)


    convertMP3ToM4A(outputFilePathMP3, outputFilePathM4A)

    return nil
}

func convertMP3ToM4A(inputFile, outputFile string) error {
	cmd := exec.Command("ffmpeg", "-i", inputFile, "-c:a", "aac", outputFile)
	return cmd.Run()
}

func detectLanguage(text string) (*translate.Detection, error) {
    // text := "こんにちは世界"
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