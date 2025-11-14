package graph_test

import (
    "context"
    "testing"

    "github.com/beeline/repodoc/internal/ast"
    "github.com/beeline/repodoc/internal/graph"
    "github.com/beeline/repodoc/internal/git"
    "github.com/stretchr/testify/require"
)

func TestBuilderCreatesNodesAndEdges(t *testing.T) {
    astReport := &ast.Report{
        Files: []string{"a.go"},
        Imports: map[string][]string{"a.go": {`"fmt"`}},
        Calls: map[string][]string{"a.go": {"fmt.Println"}},
    }
    gitReport := &git.Report{
        CommitRanges: map[string][]string{"a.go": {"c1"}},
    }
    builder := graph.NewBuilder()
    g, err := builder.Build(context.Background(), astReport, gitReport)
    require.NoError(t, err)
    require.Len(t, g.Nodes, 1)
    require.Len(t, g.Edges, 2)
    require.Contains(t, g.Adj["a.go"], `fmt.Println`)
}
