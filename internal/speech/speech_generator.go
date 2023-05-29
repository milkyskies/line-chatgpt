package speech

import (
	"context"
	"fmt"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"cloud.google.com/go/translate"
)

type SpeechGenerator struct {
	client       *texttospeech.Client
	language     *translate.Detection
	FileManager  *FileManager
	TextToSpeech texttospeechpb.SynthesizeSpeechRequest
}

func NewSpeechGenerator(fm *FileManager) *SpeechGenerator {
	return &SpeechGenerator{FileManager: fm}
}

func (sg *SpeechGenerator) Init(message string) error {
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return err
	}
	sg.client = client

	language, err := sg.detectLanguage(message)
	if err != nil {
		return err
	}
	sg.language = language

	err = sg.FileManager.MakeOutputDir()
	if err != nil {
		return err
	}

	return nil
}

func (sg *SpeechGenerator) GenerateAudio(message string, id string) error {
	voiceSelectionParams := sg.getVoiceSelectionParams()

	sg.TextToSpeech = sg.buildSynthesisRequest(message, &voiceSelectionParams)

	resp, err := sg.client.SynthesizeSpeech(context.Background(), &sg.TextToSpeech)
	if err != nil {
		return fmt.Errorf("could not synthesize speech: %f", err)
	}

	err = sg.FileManager.MakeOutputDir()
	if err != nil {
		return err
	}

	outputFilePathMP3, err := sg.FileManager.WriteAudioToFile(id, resp.AudioContent)
	if err != nil {
		return err
	}

	err = sg.FileManager.ConvertMP3ToM4A(outputFilePathMP3, id)

	return err
}

func (sg *SpeechGenerator) getVoiceSelectionParams() texttospeechpb.VoiceSelectionParams {
	langTag := sg.language.Language.String()

	switch langTag {
	case "ja":
		return texttospeechpb.VoiceSelectionParams{
			LanguageCode: "ja-JP",
			Name:         "ja-JP-Neural2-C",
		}
	case "en":
		return texttospeechpb.VoiceSelectionParams{
			LanguageCode: "en-US",
			Name:         "en-US-Standard-A",
		}
	default:
		return texttospeechpb.VoiceSelectionParams{
			LanguageCode: langTag,
			SsmlGender:   texttospeechpb.SsmlVoiceGender_MALE,
		}
	}
}

func (sg *SpeechGenerator) buildSynthesisRequest(message string, voice *texttospeechpb.VoiceSelectionParams) texttospeechpb.SynthesizeSpeechRequest {
	return texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: message},
		},
		Voice: voice,
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3},
	}
}

func (sg *SpeechGenerator) detectLanguage(text string) (*translate.Detection, error) {
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
