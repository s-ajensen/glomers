package broadcast

import (
	"encoding/json"
	"glomers/src/errors"
	"slices"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

const initialSize = 32

type Broadcast struct {
	mu     sync.RWMutex
	msgSet map[int]struct{}
	msgs   []int
}

func NewBroadcast() *Broadcast {
	return &Broadcast{msgSet: make(map[int]struct{}, initialSize), msgs: make([]int, 0, initialSize)}
}

func (b *Broadcast) Handlers(n *maelstrom.Node) map[string]maelstrom.HandlerFunc {
	return map[string]maelstrom.HandlerFunc{
		"broadcast": BroadcastHandler(n, b),
		"read":      ReadHandler(n, b),
	}
}

func (b *Broadcast) Messages() []int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	msgs := slices.Clone(b.msgs)
	return msgs
}

func (b *Broadcast) add(msg int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, seen := b.msgSet[msg]; seen {
		return
	}
	b.msgSet[msg] = struct{}{}
	b.msgs = append(b.msgs, msg)
}

type BroadcastReq struct {
	Message int `json:"message"`
}

type BroadcastReply struct {
	Type string `json:"type"`
}

func BroadcastHandler(node *maelstrom.Node, b *Broadcast) maelstrom.HandlerFunc {
	return func(msg maelstrom.Message) error {
		var req BroadcastReq
		if err := json.Unmarshal(msg.Body, &req); err != nil {
			return node.Reply(msg, errors.Malformed(err.Error()))
		}

		b.add(req.Message)

		return node.Reply(msg, BroadcastReply{Type: "broadcast_ok"})
	}
}

type ReadReply struct {
	Type     string `json:"type"`
	Messages []int  `json:"messages"`
}

func ReadHandler(node *maelstrom.Node, b *Broadcast) maelstrom.HandlerFunc {
	return func(msg maelstrom.Message) error {
		return node.Reply(msg, ReadReply{Type: "read_ok", Messages: b.Messages()})
	}
}
