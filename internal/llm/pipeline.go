package llm

import (
    "context"
    "fmt"
    "sync"

    "github.com/beeline/repodoc/internal/docs"
)

type LLMClient interface {
    Call(context.Context, string, string) (string, error)
}

type Pipeline struct {
    client LLMClient
}

func NewPipeline(client LLMClient) *Pipeline {
    return &Pipeline{client: client}
}

func (p *Pipeline) Generate(ctx context.Context, repoID string, docParts []docs.Document) error {
    var wg sync.WaitGroup
    sem := make(chan struct{}, 3)
    for _, part := range docParts {
        wg.Add(1)
        sem <- struct{}{}
        go func(part docs.Document) {
            defer wg.Done()
            defer func() { <-sem }()
            resp, err := p.client.Call(ctx, part.Title, part.Content)
            if err != nil {
                fmt.Printf("llm call failed: %v\n", err)
                return
            }
            if err := part.Store(resp); err != nil {
                fmt.Printf("store doc failed: %v\n", err)
            }
        }(part)
    }
    wg.Wait()
    return nil
}
