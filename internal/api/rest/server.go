package rest

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "sync"

    "github.com/beeline/repodoc/internal/ast"
    "github.com/beeline/repodoc/internal/graph"
    "github.com/beeline/repodoc/internal/ingestion"
    "github.com/gin-gonic/gin"
)

type Server struct {
    router *gin.Engine
    ingestion *ingestion.Service
    graphCache *sync.Map
    astAnalyzer *ast.Analyzer
}

func NewServer(ingestion *ingestion.Service, astAnalyzer *ast.Analyzer, graphStore *sync.Map) *Server {
    router := gin.New()
    router.Use(gin.Recovery())
    server := &Server{router: router, ingestion: ingestion, astAnalyzer: astAnalyzer, graphCache: graphStore}
    router.POST("/repos/upload", server.uploadRepo)
    router.POST("/repos/:id/analyze", server.triggerAnalysis)
    router.GET("/repos/:id/status", server.status)
    router.GET("/repos/:id/graph", server.graph)
    router.GET("/repos/:id/docs", server.listDocs)
    router.GET("/repos/:id/docs/:name", server.fetchDoc)
    router.GET("/repos/:id/ast/*file", server.astView)
    return server
}

func (s *Server) uploadRepo(c *gin.Context) {
    var payload struct {
        URL string `json:"url" binding:"required"`
    }
    if err := c.ShouldBindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    id, err := s.ingestion.RegisterRepository(c.Request.Context(), payload.URL)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusAccepted, gin.H{"repo_id": id})
}

func (s *Server) triggerAnalysis(c *gin.Context) {
    repoID := c.Param("id")
    var payload struct {
        URL string `json:"url" binding:"required"`
    }
    if err := c.ShouldBindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    graphData, err := s.ingestion.AnalyzeRepository(c.Request.Context(), repoID, payload.URL)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    s.graphCache.Store(repoID, graphData)
    c.JSON(http.StatusOK, gin.H{"status": "analysis_started"})
}

func (s *Server) status(c *gin.Context) {
    repoID := c.Param("id")
    if _, ok := s.graphCache.Load(repoID); ok {
        c.JSON(http.StatusOK, gin.H{"status": "ready"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "processing"})
}

func (s *Server) graph(c *gin.Context) {
    repoID := c.Param("id")
    val, ok := s.graphCache.Load(repoID)
    if !ok {
        c.JSON(http.StatusNotFound, gin.H{"error": "graph not available"})
        return
    }
    graphData := val.(*graph.Graph)
    c.JSON(http.StatusOK, gin.H{"nodes": graphData.Nodes, "edges": graphData.Edges, "adj": graphData.Adj})
}

func (s *Server) listDocs(c *gin.Context) {
    repoID := c.Param("id")
    dir := filepath.Join("docs", repoID)
    files, err := os.ReadDir(dir)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"docs": []string{}})
        return
    }
    names := make([]string, 0, len(files))
    for _, f := range files {
        names = append(names, f.Name())
    }
    c.JSON(http.StatusOK, gin.H{"docs": names})
}

func (s *Server) fetchDoc(c *gin.Context) {
    repoID := c.Param("id")
    name := c.Param("name")
    path := filepath.Join("docs", repoID, fmt.Sprintf("%s.md", name))
    c.File(path)
}

func (s *Server) astView(c *gin.Context) {
    repoID := c.Param("id")
    filePath := c.Param("file")
    absPath := filepath.Join("/tmp", repoID, filePath)
    report, err := s.astAnalyzer.Run(c.Request.Context(), absPath)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, report)
}

func (s *Server) Run(ctx context.Context, addr string) error {
    return s.router.Run(addr)
}
