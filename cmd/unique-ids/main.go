package main

import (
	"log"

	generate "glomers/src/unique_ids"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	node := maelstrom.NewNode()
	node.Handle("generate", generate.Handler(node))
	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
