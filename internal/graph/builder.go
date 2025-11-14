package graph

import (
    "context"
    "strings"

    "github.com/beeline/repodoc/internal/ast"
    "github.com/beeline/repodoc/internal/git"
)

type Node struct {
    ID string
    Type string
    Label string
}

type Edge struct {
    From string
    To string
    Type string
}

type Graph struct {
    Nodes []Node
    Edges []Edge
    Adj   map[string][]string
    Dot   string
}

type Builder struct{}

func NewBuilder() *Builder {
    return &Builder{}
}

func (b *Builder) Build(ctx context.Context, astReport *ast.Report, gitReport *git.Report) (*Graph, error) {
    nodes := map[string]Node{}
    edges := []Edge{}
    adj := map[string][]string{}
    for _, file := range astReport.Files {
        nodes[file] = Node{ID: file, Type: "file", Label: file}
    }
    for file, imports := range astReport.Imports {
        for _, imp := range imports {
            from := file
            to := strings.Trim(imp, "\"")
            edges = append(edges, Edge{From: from, To: to, Type: "import"})
            adj[from] = append(adj[from], to)
        }
    }
    for file, calls := range astReport.Calls {
        for _, call := range calls {
            edges = append(edges, Edge{From: file, To: call, Type: "call"})
            adj[file] = append(adj[file], call)
        }
    }
    for file := range gitReport.CommitRanges {
        if _, ok := nodes[file]; !ok {
            nodes[file] = Node{ID: file, Type: "file", Label: file}
        }
    }
    dot := b.toDot(nodes, edges)
    graph := &Graph{
        Nodes: mapValues(nodes),
        Edges: edges,
        Adj:   adj,
        Dot:   dot,
    }
    return graph, nil
}

func (b *Builder) toDot(nodes map[string]Node, edges []Edge) string {
    sb := strings.Builder{}
    sb.WriteString("digraph repodoc {\n")
    for _, node := range nodes {
        sb.WriteString("  \"")
        sb.WriteString(node.ID)
        sb.WriteString("\" [label=\"")
        sb.WriteString(node.Label)
        sb.WriteString("\", shape=box];\n")
    }
    for _, edge := range edges {
        sb.WriteString("  \"")
        sb.WriteString(edge.From)
        sb.WriteString("\" -> \"")
        sb.WriteString(edge.To)
        sb.WriteString("\" [label=\"")
        sb.WriteString(edge.Type)
        sb.WriteString("\"];\n")
    }
    sb.WriteString("}\n")
    return sb.String()
}

func mapValues[M ~map[K]V, K comparable, V any](m M) []V {
    vals := make([]V, 0, len(m))
    for _, v := range m {
        vals = append(vals, v)
    }
    return vals
}
