package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/VeneLooool/drone-client/internal/config"
	"github.com/VeneLooool/drone-client/internal/model"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

const (
	externalDroneEventsTopic = "external-drone-events"
	retryCount               = 3
	retryDelay               = 500 * time.Millisecond
)

type Publisher struct {
	writer writer
}

func New(ctx context.Context, cfg *config.KafkaConfig) *Publisher {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(fmt.Sprintf("%s:%s", cfg.KafkaHost, cfg.KafkaPort)),
		Topic:                  externalDroneEventsTopic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	return &Publisher{
		writer: w,
	}
}

func (p *Publisher) Publish(ctx context.Context, event model.Event) (err error) {
	marshalledEvent, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Key:   []byte(event.GetEventKey()),
		Value: marshalledEvent,
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	for range retryCount {
		err := p.writer.WriteMessages(ctx, message)
		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(retryDelay)
			continue
		}
		if err != nil {
			log.Printf("failed to write messages: %s", err.Error())
		}
		break
	}
	return nil
}

func (p *Publisher) Close() error {
	return p.writer.Close()
}
