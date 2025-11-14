package git_test

import (
    "context"
    "path/filepath"
    "testing"

    "github.com/beeline/repodoc/internal/git"
)

func TestAnalyzerProcessesCommits(t *testing.T) {
    analyzer := git.NewAnalyzer()
    repoPath := filepath.Join("..", "..", "testdata", "gitrepo")
    report, err := analyzer.Run(context.Background(), repoPath)
    if err != nil {
        t.Fatalf("git analyzer error: %v", err)
    }
    if report.Commits == 0 {
        t.Fatalf("expected commits")
    }
}
