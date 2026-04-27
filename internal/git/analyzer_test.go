package git_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/beeline/repodoc/internal/git"
)

func TestAnalyzerProcessesCommits(t *testing.T) {
	repoPath := t.TempDir()

	repo, err := gogit.PlainInit(repoPath, false)
	if err != nil {
		t.Fatalf("init repo: %v", err)
	}

	filePath := filepath.Join(repoPath, "main.go")
	if err := os.WriteFile(filePath, []byte("package main\n"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	wt, err := repo.Worktree()
	if err != nil {
		t.Fatalf("worktree: %v", err)
	}
	if _, err := wt.Add("main.go"); err != nil {
		t.Fatalf("add file: %v", err)
	}

	_, err = wt.Commit("init", &gogit.CommitOptions{
		Author: &object.Signature{
			Name:  "Tester",
			Email: "tester@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		t.Fatalf("commit: %v", err)
	}

	analyzer := git.NewAnalyzer()
	report, err := analyzer.Run(context.Background(), repoPath)
	if err != nil {
		t.Fatalf("git analyzer error: %v", err)
	}
	if report.Commits == 0 {
		t.Fatalf("expected commits")
	}
	if report.Authors["tester@example.com"] == 0 {
		t.Fatalf("expected author stats")
	}
}
