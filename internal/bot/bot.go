package bot

type Bot interface {
    GenerateReply(input string) (string, error)
}

