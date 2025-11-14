package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "sync"

    "github.com/beeline/repodoc/configs"
    grpcapi "github.com/beeline/repodoc/internal/api/grpc"
    "github.com/beeline/repodoc/internal/api/rest"
    "github.com/beeline/repodoc/internal/ast"
    kafkaevents "github.com/beeline/repodoc/internal/events/kafka"
    "github.com/beeline/repodoc/internal/git"
    "github.com/beeline/repodoc/internal/graph"
    "github.com/beeline/repodoc/internal/ingestion"
    "github.com/beeline/repodoc/internal/llm"
    "github.com/beeline/repodoc/internal/docs"
    "github.com/beeline/repodoc/internal/storage/clickhouse"
    "github.com/beeline/repodoc/internal/storage/minio"
    "github.com/beeline/repodoc/internal/storage/postgres"
    redisclient "github.com/beeline/repodoc/internal/storage/redis"
    "github.com/beeline/repodoc/pkg/telemetry"
    "google.golang.org/grpc"
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

    _, err = clickhouse.New(ctx, cfg.Storage)
    if err != nil {
        log.Fatalf("clickhouse: %v", err)
    }
    _, err = redisclient.New(ctx, cfg.Storage)
    if err != nil {
        log.Fatalf("redis: %v", err)
    }
    _, err = minio.New(ctx, cfg.Storage)
    if err != nil {
        log.Fatalf("minio: %v", err)
    }

    astAnalyzer := ast.NewAnalyzer()
    gitAnalyzer := git.NewAnalyzer()
    graphBuilder := graph.NewBuilder()
    yandex := llm.NewYandexClient(cfg.LLM.YandexAPIURL, cfg.LLM.YandexAPIKey)
    llmPipeline := llm.NewPipeline(yandex)
    docsGenerator := docs.NewGenerator()
    kafkaProducer := kafkaevents.NewProducer(cfg.Event.KafkaBrokers, fmt.Sprintf("%s.events", cfg.Event.TopicPrefix))
    ingestionService := ingestion.NewService(pg, kafkaProducer, astAnalyzer, gitAnalyzer, graphBuilder, llmPipeline, docsGenerator)

    cache := &sync.Map{}
    restServer := rest.NewServer(ingestionService, astAnalyzer, cache)
    grpcServer := grpc.NewServer()
    grpcSvc := grpcapi.NewServer(ingestionService, astAnalyzer, cache)
    grpcSvc.Register(grpcServer)

    lis, err := net.Listen("tcp", cfg.GrpcAddress())
    if err != nil {
        log.Fatalf("grpc listen: %v", err)
    }
    go func() {
        if err := grpcServer.Serve(lis); err != nil {
            log.Fatalf("grpc: %v", err)
        }
    }()

    if err := restServer.Run(ctx, cfg.RestAddress()); err != nil {
        log.Fatalf("rest: %v", err)
    }
}
