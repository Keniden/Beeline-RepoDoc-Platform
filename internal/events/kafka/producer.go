package kafkaevents

import (
    "context"
    "encoding/json"

    "github.com/segmentio/kafka-go"
)

type Event struct {
    Type string `json:"type"`
    RepoID string `json:"repo_id"`
    Payload map[string]any `json:"payload"`
}

type Producer struct {
    writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
    writer := kafka.NewWriter(kafka.WriterConfig{
        Brokers: brokers,
        Topic: topic,
    })
    return &Producer{writer: writer}
}

func (p *Producer) Send(ctx context.Context, event Event) error {
    payload, err := json.Marshal(event)
    if err != nil {
        return err
    }
    return p.writer.WriteMessages(ctx, kafka.Message{Value: payload})
}

func (p *Producer) Close() error {
    return p.writer.Close()
}
