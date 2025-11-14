package docs

import (
    "fmt"
    "os"
    "path/filepath"
)

type Document struct {
    RepoID string
    Title string
    Content string
    OutputDir string
}

func (d Document) Store(content string) error {
    if d.OutputDir == "" {
        d.OutputDir = filepath.Join("docs", d.RepoID)
    }
    if err := os.MkdirAll(d.OutputDir, 0o755); err != nil {
        return err
    }
    filename := filepath.Join(d.OutputDir, fmt.Sprintf("%s.md", d.Title))
    return os.WriteFile(filename, []byte(content), 0o644)
}

func (d Document) Describe() string {
    return fmt.Sprintf("document %s for %s", d.Title, d.RepoID)
}

func (d Document) Metadata() map[string]any {
    return map[string]any{
        "repo": d.RepoID,
        "title": d.Title,
    }
}
