package ast_test

import (
    "context"
    "path/filepath"
    "testing"

    "github.com/beeline/repodoc/internal/ast"
)

func TestAnalyzerRuns(t *testing.T) {
    analyzer := ast.NewAnalyzer()
    repoPath := filepath.Join("..", "..", "testdata", "astrepo")
    report, err := analyzer.Run(context.Background(), repoPath)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(report.Files) == 0 {
        t.Fatalf("expected file list")
    }
}
