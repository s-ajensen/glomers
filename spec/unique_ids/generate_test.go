package generate_test

import (
	"encoding/json"
	helpers "glomers/spec/helpers"
	generate "glomers/src/unique_ids"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func gen(n helpers.Node, client string) generate.GenerateReply {
	return n.Send[generate.GenerateReply](client, json.RawMessage(`{"type":"generate"}`))
}

func TestHandlerReturnsFirstId(t *testing.T) {
	node := helpers.Start(t, "n1", generate.Handler)
	assert.NotZero(t, gen(node, "c1").ID)
}

func TestHandlerReturnsDifferentId(t *testing.T) {
	node := helpers.Start(t, "n1", generate.Handler)
	assert.NotEqual(t, gen(node, "c1").ID, gen(node, "c1").ID)
}

func TestHandlerDifferentNodes_ReturnDifferentId(t *testing.T) {
	node1 := helpers.Start(t, "n1", generate.Handler)
	node2 := helpers.Start(t, "n2", generate.Handler)

	assert.NotEqual(t, gen(node1, "c1").ID, gen(node2, "c1").ID)
}
