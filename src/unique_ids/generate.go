package generate

import (
	"encoding/json"
	"fmt"
	"sync/atomic"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type GenerateReply struct {
	Type string           `json:"type"`
	ID   *json.RawMessage `json:"id"`
}

func Handler(node *maelstrom.Node) maelstrom.HandlerFunc {
	var next atomic.Int32
	return func(msg maelstrom.Message) error {
		id := next.Add(1)
		body := GenerateReply{
			Type: "generate_ok",
			ID:   new(json.RawMessage(fmt.Sprintf(`"%s%d"`, node.ID(), id))),
		}

		return node.Reply(msg, body)
	}
}
