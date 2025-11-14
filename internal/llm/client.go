package llm

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type YandexClient struct {
    url string
    apiKey string
    http *http.Client
}

func NewYandexClient(url, key string) *YandexClient {
    return &YandexClient{
        url: url,
        apiKey: key,
        http: &http.Client{Timeout: 30 * time.Second},
    }
}

type responsePayload struct {
    Response string `json:"response"`
}

func (c *YandexClient) Call(ctx context.Context, title, payload string) (string, error) {
    body, err := json.Marshal(map[string]string{
        "title": title,
        "input": payload,
    })
    if err != nil {
        return "", err
    }
    req, err := http.NewRequestWithContext(ctx, "POST", c.url, bytes.NewReader(body))
    if err != nil {
        return "", err
    }
    req.Header.Set("Authorization", fmt.Sprintf("Api-Key %s", c.apiKey))
    req.Header.Set("Content-Type", "application/json")
    resp, err := c.http.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("yandex gpt status %d", resp.StatusCode)
    }
    var payloadResp responsePayload
    if err := json.NewDecoder(resp.Body).Decode(&payloadResp); err != nil {
        return "", err
    }
    return payloadResp.Response, nil
}
