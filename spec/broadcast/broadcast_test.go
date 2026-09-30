package broadcast_test

import (
	"encoding/json"
	"fmt"
	"glomers/spec/helpers"
	"glomers/src/broadcast"
	"sync"
	"testing"

	"github.com/alecthomas/assert/v2"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func sendBroadcast(n helpers.Node, client string, message int) broadcast.BroadcastReply {
	msg := fmt.Sprintf(`{"type":"broadcast","message":%d}`, message)
	return n.Send[broadcast.BroadcastReply](client, json.RawMessage(msg))
}

func sendRead(n helpers.Node, client string) broadcast.ReadReply {
	return n.Send[broadcast.ReadReply](client, json.RawMessage(`{"type":"read"}`))
}

func TestBroadcastHandler_RejectsMalformed(t *testing.T) {
	b := broadcast.NewBroadcast()
	node := helpers.Start(t, "n1", b)

	resp := node.Send[broadcast.BroadcastReply]("c1", json.RawMessage(`{"type":"broadcast","message":"foo"}`))
	assert.Equal(t, "error", resp.Type)
}

func TestBroadcastHandlerReturnsOK(t *testing.T) {
	b := broadcast.NewBroadcast()
	node := helpers.Start(t, "n1", b)
	assert.Equal(t, "broadcast_ok", sendBroadcast(node, "c1", 1).Type)
}

func TestBroadcastHandler_SchedulesConcurrentWrites(t *testing.T) {
	msgs := 1000
	b := broadcast.NewBroadcast()
	node := helpers.Start(t, "n1", b)

	var wg sync.WaitGroup
	handler := node.Handlers["broadcast"]
	for i := range msgs {
		wg.Go(func() {
			body := json.RawMessage(fmt.Sprintf(`{"type":"broadcast","message":%d}`, i))
			handler(maelstrom.Message{Body: body})
		})
	}
	wg.Wait()

	assert.Equal(t, msgs, len(b.Messages()))
}

func TestReadHandlerReturnsNoMessages(t *testing.T) {
	b := broadcast.NewBroadcast()
	node := helpers.Start(t, "n1", b)
	resp := sendRead(node, "c1")
	assert.True(t, len(resp.Messages) == 0)
}

func TestReadHandlerReturnsMessage(t *testing.T) {
	b := broadcast.NewBroadcast()
	node := helpers.Start(t, "n1", b)
	sendBroadcast(node, "c1", 1)
	resp := sendRead(node, "c1")
	assert.SliceContains(t, resp.Messages, 1)
}
