package kafkaevents

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/segmentio/kafka-go"
)

type Consumer struct {
    reader *kafka.Reader
    handler func(context.Context, Event) error
}

func NewConsumer(brokers []string, topic, groupID string, handler func(context.Context, Event) error) *Consumer {
    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers: brokers,
        Topic: topic,
        GroupID: groupID,
    })
    return &Consumer{reader: reader, handler: handler}
}

func (c *Consumer) Start(ctx context.Context) error {
    for {
        msg, err := c.reader.FetchMessage(ctx)
        if err != nil {
            return fmt.Errorf("fetch event: %w", err)
        }
        var evt Event
        if err := json.Unmarshal(msg.Value, &evt); err != nil {
            _ = c.reader.CommitMessages(ctx, msg)
            continue
        }
        if err := c.handler(ctx, evt); err != nil {
            fmt.Printf("kafka handler error: %v\n", err)
        }
        if err := c.reader.CommitMessages(ctx, msg); err != nil {
            return err
        }
    }
}

func (c *Consumer) Close() error {
    return c.reader.Close()
}
