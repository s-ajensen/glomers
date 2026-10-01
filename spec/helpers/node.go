package helpers

import (
	"bytes"
	"encoding/json"
	"glomers/src/workload"
	"sync"
	"testing"

	"github.com/alecthomas/assert/v2"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

// provides a hook for handling messages synchronously &
// parsing stdout
type Node struct {
	t        *testing.T
	id       string
	handlers map[string]maelstrom.HandlerFunc
	Out      *json.Decoder
}

func (n Node) Send(client string, body json.RawMessage) {
	n.t.Helper()

	req := maelstrom.Message{Src: client, Dest: n.id, Body: body}
	assert.NoError(n.t, n.handlers[msgType(n.t, body)](req))
}

func (n Node) ReadReply[T any]() T {
	var msg maelstrom.Message
	assert.NoError(n.t, n.Out.Decode(&msg))
	var reply T
	assert.NoError(n.t, json.Unmarshal(msg.Body, &reply))
	return reply
}

type MsgFn func(idx int) json.RawMessage

func (n Node) SendConcurrently[T any](routines int, msgFn MsgFn) []maelstrom.Message {
	n.t.Helper()
	var wg sync.WaitGroup
	for i := range routines {
		wg.Go(func() {
			body := msgFn(i)
			msg := maelstrom.Message{Body: body}
			assert.NoError(n.t, n.handlers[msgType(n.t, body)](msg))
		})
	}
	wg.Wait()

	msgs := make([]maelstrom.Message, routines)
	for i := range msgs {
		assert.NoError(n.t, n.Out.Decode(&msgs[i]))
	}
	return msgs
}

func msgType(t *testing.T, body json.RawMessage) string {
	var reqType struct {
		Type string `json:"type"`
	}
	assert.NoError(t, json.Unmarshal(body, &reqType))
	return reqType.Type
}

func Start(t *testing.T, id string, peers []string, workload workload.Workload) Node {
	t.Helper()
	var out bytes.Buffer
	n := maelstrom.NewNode()
	n.Stdout = &out
	n.Init(id, peers)

	return Node{t: t, id: id, handlers: workload.Handlers(n), Out: json.NewDecoder(&out)}
}
