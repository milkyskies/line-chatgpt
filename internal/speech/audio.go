package speech

import (
	"fmt"
	"os/exec"
)

func convertMP3ToM4A(inputFile, outputFile string) error {
	cmd := exec.Command("ffmpeg", "-i", inputFile, "-c:a", "aac", outputFile)

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
