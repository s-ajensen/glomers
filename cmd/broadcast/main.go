package main

import (
	"glomers/src/broadcast"
	"glomers/src/workload"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	workload.Start(maelstrom.NewNode(), broadcast.NewBroadcast())
}
