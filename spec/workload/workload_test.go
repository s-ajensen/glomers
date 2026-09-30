package workload_test

import (
	"bytes"
	"encoding/json"
	"glomers/src/workload"
	"strings"
	"testing"

	"github.com/alecthomas/assert/v2"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Foo struct{}

func (f Foo) Handlers(n *maelstrom.Node) map[string]maelstrom.HandlerFunc {
	return map[string]maelstrom.HandlerFunc{
		"foo": func(msg maelstrom.Message) error {
			return n.Reply(msg, json.RawMessage(`{"bar":"baz"}`))
		},
	}
}

func TestWorkloadConfiguresHandlers(t *testing.T) {
	f := Foo{}

	msg := maelstrom.Message{Src: "c1", Dest: "n1", Body: json.RawMessage(`{"type":"foo"}`)}
	msgBuf, err := json.Marshal(msg)
	assert.NoError(t, err)

	in := strings.NewReader(string(msgBuf))
	var out bytes.Buffer
	node := maelstrom.NewNode()
	node.Stdin = in
	node.Stdout = &out
	node.Init("n1", []string{"n1"})

	workload.Start(node, f)

	var resp maelstrom.Message
	json.Unmarshal(out.Bytes(), &resp)
	assert.Contains(t, string(resp.Body), `"bar":"baz"`)
}
