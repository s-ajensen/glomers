package main

import (
	"log"

	"glomers/src/broadcast"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	var messages [1]int
	node := maelstrom.NewNode()
	node.Handle("broadcast", broadcast.BroadcastHandler(node, messages[:]))
	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
