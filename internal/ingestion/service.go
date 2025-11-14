package ingestion

import (
    "context"
    "fmt"
    "path/filepath"

    "github.com/google/uuid"
    gogit "github.com/go-git/go-git/v5"
    "github.com/beeline/repodoc/internal/events/kafka"
    "github.com/beeline/repodoc/internal/git"
    "github.com/beeline/repodoc/internal/ast"
    "github.com/beeline/repodoc/internal/graph"
    "github.com/beeline/repodoc/internal/llm"
    "github.com/beeline/repodoc/internal/docs"
    "github.com/beeline/repodoc/internal/storage/postgres"
)

type Service struct {
    pg          *postgres.Store
    kafka       *kafkaevents.Producer
    astAnalyzer *ast.Analyzer
    gitAnalyzer *git.Analyzer
    graphBuilder *graph.Builder
    llmPipeline *llm.Pipeline
    docsGen *docs.Generator
}

func NewService(pg *postgres.Store, kafka *kafkaevents.Producer, astAnalyzer *ast.Analyzer, gitAnalyzer *git.Analyzer, graphBuilder *graph.Builder, llmPipeline *llm.Pipeline, docsGen *docs.Generator) *Service {
    return &Service{pg: pg, kafka: kafka, astAnalyzer: astAnalyzer, gitAnalyzer: gitAnalyzer, graphBuilder: graphBuilder, llmPipeline: llmPipeline, docsGen: docsGen}
}

func (s *Service) RegisterRepository(ctx context.Context, repoURL string) (string, error) {
    id := uuid.NewString()
    if err := s.pg.SaveRepo(ctx, repoURL); err != nil {
        return "", err
    }
    return id, nil
}

func (s *Service) AnalyzeRepository(ctx context.Context, repoID, repoURL string) (*graph.Graph, error) {
    repoPath := filepath.Join("/tmp", repoID)
    _, err := gogit.PlainCloneContext(ctx, repoPath, false, &gogit.CloneOptions{URL: repoURL, Depth: 1})
    if err != nil {
        return nil, fmt.Errorf("clone repo: %w", err)
    }

    astReport, err := s.astAnalyzer.Run(ctx, repoPath)
    if err != nil {
        return nil, err
    }

    gitReport, err := s.gitAnalyzer.Run(ctx, repoPath)
    if err != nil {
        return nil, err
    }

    graphData, err := s.graphBuilder.Build(ctx, astReport, gitReport)
    if err != nil {
        return nil, err
    }

    if err := s.kafka.Send(ctx, kafkaevents.Event{Type: "module_completed", RepoID: repoID, Payload: map[string]any{"graph": graphData}}); err != nil {
        return nil, err
    }

    docs, err := s.docsGen.Generate(ctx, repoID, astReport, gitReport, graphData)
    if err != nil {
        return nil, err
    }

    if err := s.llmPipeline.Generate(ctx, repoID, docs); err != nil {
        return nil, err
    }

    return graphData, nil
}
