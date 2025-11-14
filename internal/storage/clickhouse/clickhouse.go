package clickhouse

import (
    "context"
    "fmt"

    "github.com/ClickHouse/clickhouse-go/v2"
    "github.com/beeline/repodoc/configs"
)

type Client struct {
    conn clickhouse.Conn
}

func New(ctx context.Context, cfg configs.StorageConfig) (*Client, error) {
    conn, err := clickhouse.Open(&clickhouse.Options{
        Addr: []string{cfg.ClickHouseDSN},
        Auth: clickhouse.Auth{Database: "default"},
    })
    if err != nil {
        return nil, fmt.Errorf("connect clickhouse: %w", err)
    }
    return &Client{conn: conn}, nil
}

func (c *Client) Close(ctx context.Context) error {
    if c.conn != nil {
        return c.conn.Close()
    }
    return nil
}
