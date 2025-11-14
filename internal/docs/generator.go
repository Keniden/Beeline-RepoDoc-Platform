package docs

import (
    "context"
    "fmt"
    "strings"

    "github.com/beeline/repodoc/internal/ast"
    "github.com/beeline/repodoc/internal/git"
    "github.com/beeline/repodoc/internal/graph"
)

type Generator struct{
}

func NewGenerator() *Generator {
    return &Generator{}
}

func (g *Generator) Generate(ctx context.Context, repoID string, astReport *ast.Report, gitReport *git.Report, graphData *graph.Graph) ([]Document, error) {
    docs := []Document{}
    moduleSummary := strings.Builder{}
    moduleSummary.WriteString(fmt.Sprintf("# %s overview\n", repoID))
    moduleSummary.WriteString("## Files\n")
    for _, node := range graphData.Nodes {
        moduleSummary.WriteString(fmt.Sprintf("- %s (%s)\n", node.Label, node.Type))
    }
    moduleSummary.WriteString("## Commits\n")
    moduleSummary.WriteString(fmt.Sprintf("Total commits: %d\n", gitReport.Commits))
    docs = append(docs, Document{RepoID: repoID, Title: "project_overview", Content: moduleSummary.String()})
    return docs, nil
}
