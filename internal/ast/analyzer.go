package ast

import (
    "context"
    "fmt"
    "go/ast"
    "go/parser"
    "go/token"
    "os"
    "path/filepath"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
)

type Report struct {
    Files []string
    Functions map[string][]string
    Structs map[string][]string
    Interfaces map[string][]string
    Imports map[string][]string
    Calls map[string][]string
    Cohesion map[string]float64
    Coupling map[string]float64
}

type Analyzer struct {
    tracerCtx attribute.KeyValue
}

func NewAnalyzer() *Analyzer {
    return &Analyzer{}
}

func (a *Analyzer) Run(ctx context.Context, root string) (*Report, error) {
    tracer := otel.Tracer("ast-analyzer")
    ctx, span := tracer.Start(ctx, "ast-run")
    defer span.End()

    report := &Report{
        Functions:  make(map[string][]string),
        Structs:    make(map[string][]string),
        Interfaces: make(map[string][]string),
        Imports:    make(map[string][]string),
        Calls:      make(map[string][]string),
        Cohesion:   make(map[string]float64),
        Coupling:   make(map[string]float64),
    }
    fset := token.NewFileSet()
    walk := func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() || filepath.Ext(path) != ".go" {
            return nil
        }
        report.Files = append(report.Files, path)
        fileNode, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
        if err != nil {
            span.SetAttributes(attribute.String("parser.error", err.Error()))
            return err
        }
        ast.Inspect(fileNode, func(node ast.Node) bool {
            switch n := node.(type) {
            case *ast.FuncDecl:
                report.Functions[path] = append(report.Functions[path], n.Name.Name)
                for _, stmt := range n.Body.List {
                    if call, ok := stmt.(*ast.ExprStmt); ok {
                        report.Calls[path] = append(report.Calls[path], fmt.Sprintf("%T", call))
                    }
                }
            case *ast.TypeSpec:
                switch n.Type.(type) {
                case *ast.StructType:
                    report.Structs[path] = append(report.Structs[path], n.Name.Name)
                case *ast.InterfaceType:
                    report.Interfaces[path] = append(report.Interfaces[path], n.Name.Name)
                }
            case *ast.ImportSpec:
                report.Imports[path] = append(report.Imports[path], n.Path.Value)
            }
            return true
        })
        report.Cohesion[path] = float64(len(report.Functions[path])) / float64(maxInt(1, len(report.Structs[path])+1))
        report.Coupling[path] = float64(len(report.Imports[path])) / float64(maxInt(1, len(report.Functions[path]))+1)
        return nil
    }
    if err := filepath.Walk(root, walk); err != nil {
        return nil, fmt.Errorf("walk repo: %w", err)
    }
    return report, nil
}

func maxInt(a, b int) int {
    if a > b {
        return a
    }
    return b
}
