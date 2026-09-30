package generate

import (
	"encoding/json"
	"fmt"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func ptr[T any](v T) *T { return &v }

type GenerateReply struct {
	Type string           `json:"type"`
	ID   *json.RawMessage `json:"id"`
}

func Handler(node *maelstrom.Node) maelstrom.HandlerFunc {
	next := 0
	return func(msg maelstrom.Message) error {
		next++
		body := GenerateReply{
			Type: "generate_ok",
			ID:   ptr(json.RawMessage(fmt.Sprintf(`"%s%d"`, node.ID(), next))),
		}

		return node.Reply(msg, body)
	}
}
