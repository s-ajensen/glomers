package broadcast_test

import (
	"encoding/json"
	"fmt"
	helpers "glomers/spec/helpers"
	"glomers/src/broadcast"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func sendBroadcast(n helpers.Node, client string, message int) broadcast.BroadcastReply {
	msg := fmt.Sprintf(`{"type":"broadcast","message":%d}`, message)
	return n.Send[broadcast.BroadcastReply](client, json.RawMessage(msg))
}

func broadcastHandler(messages []int) helpers.HandlerFactory {
	return func(node *maelstrom.Node) maelstrom.HandlerFunc {
		return broadcast.BroadcastHandler(node, messages)
	}
}

func sendRead(n helpers.Node, client string) broadcast.ReadReply {
	return n.Send[broadcast.ReadReply](client, json.RawMessage(`{"type":"read"}`))
}

// func TestBroadcastHandlerReturnsOK(t *testing.T) {
// 	messages := []int{}
// 	node := helpers.Start(t, "n1", broadcastHandler(messages))
// 	assert.Equal(t, "broadcast_ok", sendBroadcast(node, "c1", 1).Type)
// }

// func TestReadHandlerReturnsNoMessages(t *testing.T) {
// 	messages := []int{}
// 	node := helpers.Start(t, "n1", broadcastHandler(messages))
// 	resp := sendRead(node, "c1")
// 	assert.True(t, len(resp.Messages) == 0)
// }

// func TestReadHandlerReturnsMessage(t *testing.T) {
// 	var messages [1]int
// 	node := helpers.Start(t, "n1", broadcastHandler(messages[:]))
// 	sendBroadcast(node, "c1", 1)
// 	resp := sendRead(node, "c1")
// 	assert.SliceContains(t, resp.Messages, 1)
// }
