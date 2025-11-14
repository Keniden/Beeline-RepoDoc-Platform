package kafkaevents_test

import (
    "encoding/json"
    "testing"

    kafkaevents "github.com/beeline/repodoc/internal/events/kafka"
)

func TestEventMarshalling(t *testing.T) {
    evt := kafkaevents.Event{Type: "file_analyzed", RepoID: "r", Payload: map[string]any{"key": "value"}}
    data, err := json.Marshal(evt)
    if err != nil {
        t.Fatalf("marshal: %v", err)
    }
    var round kafkaevents.Event
    if err := json.Unmarshal(data, &round); err != nil {
        t.Fatalf("unmarshal: %v", err)
    }
    if round.Type != evt.Type || round.RepoID != evt.RepoID {
        t.Fatalf("unexpected roundtrip")
    }
}
