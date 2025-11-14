package pb

import (
    "context"

    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type Repo struct {
    Id string
    Name string
    Url string
    Status string
}

type File struct {
    Id string
    Path string
    Imports []string
}

type Edge struct {
    From string
    To string
    Type string
}

type UploadRequest struct {
    Url string
}

type UploadResponse struct {
    RepoId string
}

type AnalyzeRequest struct {
    RepoId string
    Url string
}

type AnalyzeResponse struct {
    Status string
}

type GetStatusRequest struct {
    RepoId string
}

type StatusResponse struct {
    Status string
}

type GetGraphRequest struct {
    RepoId string
}

type GraphResponse struct {
    Edges []*Edge
    Files []*File
}

type GetDocsRequest struct {
    RepoId string
}

type DocPart struct {
    Name string
    Content string
}

type DocsResponse struct {
    Docs []*DocPart
}

type DocRequest struct {
    RepoId string
    Parts []*DocPart
}

type DocResponse struct {
    Docs []*DocPart
}

type RepoManagerServer interface {
    UploadRepo(context.Context, *UploadRequest) (*UploadResponse, error)
    AnalyzeRepo(context.Context, *AnalyzeRequest) (*AnalyzeResponse, error)
    GetStatus(context.Context, *GetStatusRequest) (*StatusResponse, error)
    GetGraph(context.Context, *GetGraphRequest) (*GraphResponse, error)
    ListDocs(context.Context, *GetDocsRequest) (*DocsResponse, error)
}

type DocServiceServer interface {
    GenerateDoc(context.Context, *DocRequest) (*DocResponse, error)
}

type UnimplementedRepoManagerServer struct{}

func (UnimplementedRepoManagerServer) UploadRepo(context.Context, *UploadRequest) (*UploadResponse, error) {
    return nil, status.Errorf(codes.Unimplemented, "method UploadRepo not implemented")
}
func (UnimplementedRepoManagerServer) AnalyzeRepo(context.Context, *AnalyzeRequest) (*AnalyzeResponse, error) {
    return nil, status.Errorf(codes.Unimplemented, "method AnalyzeRepo not implemented")
}
func (UnimplementedRepoManagerServer) GetStatus(context.Context, *GetStatusRequest) (*StatusResponse, error) {
    return nil, status.Errorf(codes.Unimplemented, "method GetStatus not implemented")
}
func (UnimplementedRepoManagerServer) GetGraph(context.Context, *GetGraphRequest) (*GraphResponse, error) {
    return nil, status.Errorf(codes.Unimplemented, "method GetGraph not implemented")
}
func (UnimplementedRepoManagerServer) ListDocs(context.Context, *GetDocsRequest) (*DocsResponse, error) {
    return nil, status.Errorf(codes.Unimplemented, "method ListDocs not implemented")
}

type UnimplementedDocServiceServer struct{}

func (UnimplementedDocServiceServer) GenerateDoc(context.Context, *DocRequest) (*DocResponse, error) {
    return nil, status.Errorf(codes.Unimplemented, "method GenerateDoc not implemented")
}

func RegisterRepoManagerServer(s *grpc.Server, srv RepoManagerServer) {
    s.RegisterService(&grpc.ServiceDesc{
        ServiceName: "repodoc.RepoManager",
        HandlerType: (*RepoManagerServer)(nil),
        Methods: []grpc.MethodDesc{
            {MethodName: "UploadRepo", Handler: uploadRepoHandler},
            {MethodName: "AnalyzeRepo", Handler: analyzeRepoHandler},
            {MethodName: "GetStatus", Handler: getStatusHandler},
            {MethodName: "GetGraph", Handler: getGraphHandler},
            {MethodName: "ListDocs", Handler: listDocsHandler},
        },
    }, srv)
}

func RegisterDocServiceServer(s *grpc.Server, srv DocServiceServer) {
    s.RegisterService(&grpc.ServiceDesc{
        ServiceName: "repodoc.DocService",
        HandlerType: (*DocServiceServer)(nil),
        Methods: []grpc.MethodDesc{
            {MethodName: "GenerateDoc", Handler: docGenerateHandler},
        },
    }, srv)
}

func uploadRepoHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
    in := new(UploadRequest)
    if err := dec(in); err != nil {
        return nil, err
    }
    if interceptor == nil {
        return srv.(RepoManagerServer).UploadRepo(ctx, in)
    }
    info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/repodoc.RepoManager/UploadRepo"}
    handler := func(ctx context.Context, req interface{}) (interface{}, error) {
        return srv.(RepoManagerServer).UploadRepo(ctx, req.(*UploadRequest))
    }
    return interceptor(ctx, in, info, handler)
}

func analyzeRepoHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
    in := new(AnalyzeRequest)
    if err := dec(in); err != nil {
        return nil, err
    }
    if interceptor == nil {
        return srv.(RepoManagerServer).AnalyzeRepo(ctx, in)
    }
    info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/repodoc.RepoManager/AnalyzeRepo"}
    handler := func(ctx context.Context, req interface{}) (interface{}, error) {
        return srv.(RepoManagerServer).AnalyzeRepo(ctx, req.(*AnalyzeRequest))
    }
    return interceptor(ctx, in, info, handler)
}

func getStatusHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
    in := new(GetStatusRequest)
    if err := dec(in); err != nil {
        return nil, err
    }
    if interceptor == nil {
        return srv.(RepoManagerServer).GetStatus(ctx, in)
    }
    info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/repodoc.RepoManager/GetStatus"}
    handler := func(ctx context.Context, req interface{}) (interface{}, error) {
        return srv.(RepoManagerServer).GetStatus(ctx, req.(*GetStatusRequest))
    }
    return interceptor(ctx, in, info, handler)
}

func getGraphHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
    in := new(GetGraphRequest)
    if err := dec(in); err != nil {
        return nil, err
    }
    if interceptor == nil {
        return srv.(RepoManagerServer).GetGraph(ctx, in)
    }
    info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/repodoc.RepoManager/GetGraph"}
    handler := func(ctx context.Context, req interface{}) (interface{}, error) {
        return srv.(RepoManagerServer).GetGraph(ctx, req.(*GetGraphRequest))
    }
    return interceptor(ctx, in, info, handler)
}

func listDocsHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
    in := new(GetDocsRequest)
    if err := dec(in); err != nil {
        return nil, err
    }
    if interceptor == nil {
        return srv.(RepoManagerServer).ListDocs(ctx, in)
    }
    info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/repodoc.RepoManager/ListDocs"}
    handler := func(ctx context.Context, req interface{}) (interface{}, error) {
        return srv.(RepoManagerServer).ListDocs(ctx, req.(*GetDocsRequest))
    }
    return interceptor(ctx, in, info, handler)
}

func docGenerateHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
    in := new(DocRequest)
    if err := dec(in); err != nil {
        return nil, err
    }
    if interceptor == nil {
        return srv.(DocServiceServer).GenerateDoc(ctx, in)
    }
    info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/repodoc.DocService/GenerateDoc"}
    handler := func(ctx context.Context, req interface{}) (interface{}, error) {
        return srv.(DocServiceServer).GenerateDoc(ctx, req.(*DocRequest))
    }
    return interceptor(ctx, in, info, handler)
}
