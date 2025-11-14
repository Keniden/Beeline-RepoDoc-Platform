package llm_test

import (
    "context"
    "testing"

    "github.com/beeline/repodoc/internal/docs"
    "github.com/beeline/repodoc/internal/llm"
)

type mockLLM struct {
    calls int
}

func (m *mockLLM) Call(ctx context.Context, title, payload string) (string, error) {
    m.calls++
    return "generated", nil
}

func TestPipeline_WriteDocs(t *testing.T) {
    client := &mockLLM{}
    pipeline := llm.NewPipeline(client)
    parts := []docs.Document{{RepoID: "test", Title: "module", Content: "payload", OutputDir: "testdata/docs"}}
    if err := pipeline.Generate(context.Background(), "test", parts); err != nil {
        t.Fatalf("pipeline error: %v", err)
    }
    if client.calls != len(parts) {
        t.Fatalf("expected llm called %d times got %d", len(parts), client.calls)
    }
}
