# syntax=docker/dockerfile:1.5
FROM golang:1.22-alpine AS builder
WORKDIR /workspace
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /workspace/bin/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /workspace/bin/worker ./cmd/worker

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /workspace/bin /bin
