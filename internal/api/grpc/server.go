package grpcapi

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    "sync"

    "github.com/beeline/repodoc/internal/ast"
    "github.com/beeline/repodoc/internal/graph"
    "github.com/beeline/repodoc/internal/ingestion"
    pb "github.com/beeline/repodoc/proto"
    "google.golang.org/grpc"
)

type Server struct {
    pb.UnimplementedRepoManagerServer
    pb.UnimplementedDocServiceServer
    ingestion *ingestion.Service
    astAnalyzer *ast.Analyzer
    graphCache *sync.Map
}

func NewServer(ingestion *ingestion.Service, astAnalyzer *ast.Analyzer, cache *sync.Map) *Server {
    return &Server{ingestion: ingestion, astAnalyzer: astAnalyzer, graphCache: cache}
}

func (s *Server) Register(grpcServer *grpc.Server) {
    pb.RegisterRepoManagerServer(grpcServer, s)
    pb.RegisterDocServiceServer(grpcServer, s)
}

func (s *Server) UploadRepo(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
    id, err := s.ingestion.RegisterRepository(ctx, req.Url)
    if err != nil {
        return nil, err
    }
    return &pb.UploadResponse{RepoId: id}, nil
}

func (s *Server) AnalyzeRepo(ctx context.Context, req *pb.AnalyzeRequest) (*pb.AnalyzeResponse, error) {
    graphData, err := s.ingestion.AnalyzeRepository(ctx, req.RepoId, req.Url)
    if err != nil {
        return nil, err
    }
    s.graphCache.Store(req.RepoId, graphData)
    return &pb.AnalyzeResponse{Status: "triggered"}, nil
}

func (s *Server) GetStatus(ctx context.Context, req *pb.GetStatusRequest) (*pb.StatusResponse, error) {
    if _, ok := s.graphCache.Load(req.RepoId); ok {
        return &pb.StatusResponse{Status: "ready"}, nil
    }
    return &pb.StatusResponse{Status: "processing"}, nil
}

func (s *Server) GetGraph(ctx context.Context, req *pb.GetGraphRequest) (*pb.GraphResponse, error) {
    val, ok := s.graphCache.Load(req.RepoId)
    if !ok {
        return nil, fmt.Errorf("not ready")
    }
    graphData := val.(*graph.Graph)
    files := make([]*pb.File, 0, len(graphData.Nodes))
    for _, node := range graphData.Nodes {
        files = append(files, &pb.File{Id: node.ID, Path: node.Label})
    }
    edges := make([]*pb.Edge, 0, len(graphData.Edges))
    for _, edge := range graphData.Edges {
        edges = append(edges, &pb.Edge{From: edge.From, To: edge.To, Type: edge.Type})
    }
    return &pb.GraphResponse{Edges: edges, Files: files}, nil
}

func (s *Server) ListDocs(ctx context.Context, req *pb.GetDocsRequest) (*pb.DocsResponse, error) {
    dir := filepath.Join("docs", req.RepoId)
    entries, err := os.ReadDir(dir)
    if err != nil {
        return &pb.DocsResponse{}, nil
    }
    docs := make([]*pb.DocPart, 0, len(entries))
    for _, entry := range entries {
        content, err := os.ReadFile(filepath.Join(dir, entry.Name()))
        if err != nil {
            continue
        }
        docs = append(docs, &pb.DocPart{Name: entry.Name(), Content: string(content)})
    }
    return &pb.DocsResponse{Docs: docs}, nil
}

func (s *Server) GenerateDoc(ctx context.Context, req *pb.DocRequest) (*pb.DocResponse, error) {
    parts := make([]*pb.DocPart, 0, len(req.Parts))
    for _, part := range req.Parts {
        filename := filepath.Join("docs", req.RepoId, fmt.Sprintf("%s.md", part.Name))
        if err := os.WriteFile(filename, []byte(part.Content), 0o644); err != nil {
            return nil, err
        }
        parts = append(parts, part)
    }
    return &pb.DocResponse{Docs: parts}, nil
}
