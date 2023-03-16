package conversation

import (
	"context"
	"errors"
	"fmt"
)

type Message struct {
	ID string
	Text string
}