package main

import (
    "context"
    "fmt"
    "log"

    "github.com/beeline/repodoc/configs"
    kafkaevents "github.com/beeline/repodoc/internal/events/kafka"
    "github.com/beeline/repodoc/internal/storage/postgres"
    "github.com/beeline/repodoc/pkg/telemetry"
)

func main() {
    ctx := context.Background()
    cfg, err := configs.Load()
    if err != nil {
        log.Fatalf("config: %v", err)
    }
    tp, err := telemetry.InitTracer(ctx, cfg.Telemetry)
    if err != nil {
        log.Fatalf("telemetry: %v", err)
    }
    defer telemetry.Shutdown(ctx, tp)

    pg, err := postgres.New(ctx, cfg.Storage)
    if err != nil {
        log.Fatalf("postgres: %v", err)
    }
    defer pg.Close()

    handler := func(ctx context.Context, event kafkaevents.Event) error {
        fmt.Printf("worker event type=%s repo=%s\n", event.Type, event.RepoID)
        switch event.Type {
        case "file_analyzed":
            return nil
        case "module_completed":
            return nil
        case "doc_generated":
            return nil
        }
        return nil
    }

    consumer := kafkaevents.NewConsumer(cfg.Event.KafkaBrokers, fmt.Sprintf("%s.events", cfg.Event.TopicPrefix), "worker-group", handler)
    if err := consumer.Start(ctx); err != nil {
        log.Fatalf("consumer: %v", err)
    }
}
