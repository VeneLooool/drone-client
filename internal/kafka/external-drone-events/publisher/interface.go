package publisher

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type writer interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}
