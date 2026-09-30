package main

import (
	generate "glomers/src/unique_ids"
	"glomers/src/workload"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	workload.Start(maelstrom.NewNode(), generate.Generate{})
}
