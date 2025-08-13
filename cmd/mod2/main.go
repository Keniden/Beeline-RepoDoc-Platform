package main

import (
	"context"
	"time"
	"to-do-list/cmd/mod2/logger"
)

func main() {
	ctx := context.Background()

	ctx = logger.ContextWithTimestamp(ctx, time.Now())

	logger.Infof(ctx, "hello, %s!", "route228") 
}
