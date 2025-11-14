package docs_test

import (
    "context"
    "testing"

    "github.com/beeline/repodoc/internal/docs"
    "github.com/beeline/repodoc/internal/graph"
    "github.com/beeline/repodoc/internal/git"
    "github.com/beeline/repodoc/internal/ast"
    "github.com/stretchr/testify/require"
)

func TestGeneratorProducesOverview(t *testing.T) {
    g := &graph.Graph{
        Nodes: []graph.Node{{ID: "f", Type: "file", Label: "f.go"}},
    }
    docsGen := docs.NewGenerator()
    docsList, err := docsGen.Generate(context.Background(), "repo", &ast.Report{}, &git.Report{Commits: 3}, g)
    require.NoError(t, err)
    require.Len(t, docsList, 1)
    require.Contains(t, docsList[0].Content, "Total commits: 3")
}
