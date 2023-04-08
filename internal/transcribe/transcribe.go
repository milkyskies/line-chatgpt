package transcribe

type Transcribe interface {
	TranscribeAudioFile(filename string) (string, error)
}
