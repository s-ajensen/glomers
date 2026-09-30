package main

import (
	"glomers/src/echo"
	"glomers/src/workload"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	workload.Start(maelstrom.NewNode(), echo.Echo{})
}
