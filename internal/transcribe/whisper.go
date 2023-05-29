package transcribe

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/milkyskies/line-chatgpt/internal/openai"
	openaiApi "github.com/sashabaranov/go-openai"
)

type Whisper struct {
	OpenAI *openai.OpenAI
}

func NewWhisper(openAI *openai.OpenAI) *Whisper {
	return &Whisper{OpenAI: openAI}
}

func (w *Whisper) TranscribeAudioFile(filename string) (string, error) {
	ctx := context.Background()

	req := openaiApi.AudioRequest{
		FilePath: filepath.Join("content/line/audio", fmt.Sprintf("%s.m4a", filename)),
		Model:    openaiApi.Whisper1,
	}

	resp, err := w.OpenAI.Client.CreateTranscription(ctx, req)
	if err != nil {
		err := fmt.Errorf("transcription error: %v", err)
		return "", err
	}

	return resp.Text, nil
}

var _ Transcribe = (*Whisper)(nil)
