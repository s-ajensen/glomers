package helpers

import (
	"bytes"
	"encoding/json"
	"glomers/src/workload"
	"testing"

	"github.com/alecthomas/assert/v2"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

// provides a hook for handling messages synchronously &
// parsing stdout
type Node struct {
	T        *testing.T
	ID       string
	Handlers map[string]maelstrom.HandlerFunc
	Out      *json.Decoder
}

func (n Node) Send[T any](client string, body json.RawMessage) T {
	n.T.Helper()

	var reqType struct {
		Type string `json:"type"`
	}
	assert.NoError(n.T, json.Unmarshal(body, &reqType))

	req := maelstrom.Message{Src: client, Dest: n.ID, Body: body}
	assert.NoError(n.T, n.Handlers[reqType.Type](req))

	var msg maelstrom.Message
	assert.NoError(n.T, n.Out.Decode(&msg))
	var reply T
	assert.NoError(n.T, json.Unmarshal(msg.Body, &reply))
	return reply
}

func Start(t *testing.T, id string, workload workload.Workload) Node {
	t.Helper()
	var out bytes.Buffer
	n := maelstrom.NewNode()
	n.Stdout = &out
	n.Init(id, []string{id})

	return Node{T: t, ID: id, Handlers: workload.Handlers(n), Out: json.NewDecoder(&out)}
}
