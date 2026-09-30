package echo

import (
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Echo struct{}

func (e Echo) Handlers(n *maelstrom.Node) map[string]maelstrom.HandlerFunc {
	return map[string]maelstrom.HandlerFunc{
		"echo": Handler(n),
	}
}

func Handler(node *maelstrom.Node) maelstrom.HandlerFunc {
	return func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		body["type"] = "echo_ok"

		return node.Reply(msg, body)
	}
}
