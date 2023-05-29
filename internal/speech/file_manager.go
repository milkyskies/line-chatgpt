package speech

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hajimehoshi/go-mp3"
)

type FileManager struct {
	OutputDir string
}

func NewFileManager() *FileManager {
	return &FileManager{
		OutputDir: "content/whisper/audio_replies",
	}
}

func (fm *FileManager) MakeOutputDir() error {
	err := os.MkdirAll(fm.OutputDir, 0755)

	return err
}

func (fm *FileManager) WriteAudioToFile(id string, audioContent []byte) (string, error) {
	outputFilePathMP3 := filepath.Join(fm.OutputDir, id+".mp3")
	if err := ioutil.WriteFile(outputFilePathMP3, audioContent, 0644); err != nil {
		return "", err
	}
	fmt.Printf("Audio content written to file: %v\n", outputFilePathMP3)
	return outputFilePathMP3, nil
}

func (fm *FileManager) ConvertMP3ToM4A(mp3Path, id string) error {
	m4aPath := filepath.Join(fm.OutputDir, id+".m4a")

	cmd := exec.Command("ffmpeg", "-i", mp3Path, "-c:a", "aac", m4aPath)

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("command did not complete successfully: %w", err)
	}

	return nil
}

func (fm *FileManager) OpenAudioFile(id string) (*os.File, error) {
	f, err := os.Open(filepath.Join(fm.OutputDir, fmt.Sprintf("%s.mp3", id)))
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fm *FileManager) CalculateAudioDuration(f *os.File) (int, error) {
	d, err := mp3.NewDecoder(f)
	if err != nil {
		return 0, err
	}

	sampleSize := 4
	samples := int(d.Length()) / sampleSize
	duration := samples / d.SampleRate() * 1000

	return duration, nil
}
