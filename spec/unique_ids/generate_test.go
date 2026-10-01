package generate_test

import (
	"encoding/json"
	helpers "glomers/spec/helpers"
	generate "glomers/src/unique_ids"
	"testing"

	"github.com/alecthomas/assert/v2"
)

type repliedId struct {
	ID any `json:"id"`
}

func sendGen(n helpers.Node, client string) generate.GenerateReply {
	n.Send(client, json.RawMessage(`{"type":"generate"}`))
	return n.ReadReply[generate.GenerateReply]()
}

func TestHandlerReturnsFirstId(t *testing.T) {
	g := generate.Generate{}
	node := helpers.Start(t, "n1", nil, g)
	assert.NotZero(t, sendGen(node, "c1").ID)
}

func TestHandlerReturnsDifferentId(t *testing.T) {
	g := generate.Generate{}
	node := helpers.Start(t, "n1", nil, g)
	assert.NotEqual(t, sendGen(node, "c1").ID, sendGen(node, "c1").ID)
}

func TestHandlerReturnsDifferentIdConcurrently(t *testing.T) {
	g := generate.Generate{}
	node := helpers.Start(t, "n1", nil, g)

	makeMsg := func(idx int) json.RawMessage {
		return json.RawMessage(`{"type":"generate"}`)
	}
	msgs := node.SendConcurrently[generate.GenerateReply](5, makeMsg)

	ids := make(map[repliedId]int, len(msgs))
	var id repliedId
	for _, msg := range msgs {
		json.Unmarshal(msg.Body, &id)
		ids[id]++
	}

	for _, ct := range ids {
		assert.False(t, ct > 1)
	}
	assert.NotEqual(t, sendGen(node, "c1").ID, sendGen(node, "c1").ID)
}

func TestHandlerDifferentNodes_ReturnDifferentId(t *testing.T) {
	g := generate.Generate{}
	node1 := helpers.Start(t, "n1", nil, g)
	node2 := helpers.Start(t, "n2", nil, g)

	assert.NotEqual(t, sendGen(node1, "c1").ID, sendGen(node2, "c1").ID)
}
