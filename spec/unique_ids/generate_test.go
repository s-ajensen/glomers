package generate_test

import (
	"encoding/json"
	helpers "glomers/spec/helpers"
	generate "glomers/src/unique_ids"
	"sync"
	"testing"

	"github.com/alecthomas/assert/v2"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type repliedId struct {
	ID any `json:"id"`
}

func sendGen(n helpers.Node, client string) generate.GenerateReply {
	return n.Send[generate.GenerateReply](client, json.RawMessage(`{"type":"generate"}`))
}

func TestHandlerReturnsFirstId(t *testing.T) {
	node := helpers.Start(t, "n1", generate.Handler)
	assert.NotZero(t, sendGen(node, "c1").ID)
}

func TestHandlerReturnsDifferentId(t *testing.T) {
	node := helpers.Start(t, "n1", generate.Handler)
	assert.NotEqual(t, sendGen(node, "c1").ID, sendGen(node, "c1").ID)
}

func TestHandlerReturnsDifferentIdConcurrently(t *testing.T) {
	msgs := 1000
	node := helpers.Start(t, "n1", generate.Handler)

	var wg sync.WaitGroup
	for range msgs {
		wg.Go(func() {
			node.Handle(maelstrom.Message{Body: json.RawMessage(`{"type":"generate"}`)})
		})
	}
	wg.Wait()

	ids := make(map[repliedId]int)
	var msg maelstrom.Message
	var id repliedId
	for range msgs {
		node.Out.Decode(&msg)
		json.Unmarshal(msg.Body, &id)
		ids[id]++
	}

	for _, ct := range ids {
		assert.False(t, ct > 1)
	}
	assert.NotEqual(t, sendGen(node, "c1").ID, sendGen(node, "c1").ID)
}

func TestHandlerDifferentNodes_ReturnDifferentId(t *testing.T) {
	node1 := helpers.Start(t, "n1", generate.Handler)
	node2 := helpers.Start(t, "n2", generate.Handler)

	assert.NotEqual(t, sendGen(node1, "c1").ID, sendGen(node2, "c1").ID)
}
