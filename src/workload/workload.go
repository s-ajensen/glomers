package workload

import (
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Workload interface {
	Handlers(n *maelstrom.Node) map[string]maelstrom.HandlerFunc
}

func Start(node *maelstrom.Node, workload Workload) {
	for msgType, handler := range workload.Handlers(node) {
		node.Handle(msgType, handler)
	}
	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
