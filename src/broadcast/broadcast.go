package broadcast

import (
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type BroadcastReply struct {
	Type string
}

func BroadcastHandler(node *maelstrom.Node, messages []int) maelstrom.HandlerFunc {
	return func(msg maelstrom.Message) error {
		return node.Reply(msg, BroadcastReply{Type: "broadcast_ok"})
	}
}

type ReadReply struct {
	Type     string
	Messages []int
}

func ReadHandler(node *maelstrom.Node, messages []int) maelstrom.HandlerFunc {
	return func(msg maelstrom.Message) error {
		return node.Reply(msg, ReadReply{Type: "read_ok", Messages: messages[0:0]})
	}
}
