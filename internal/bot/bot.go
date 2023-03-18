package bot

import "github.com/milkyskies/line-chatgpt/internal/database"


type Bot interface {
    GenerateReply(input string, history []database.Message) (string, error)
}

