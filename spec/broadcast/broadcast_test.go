package broadcast_test

import (
	"encoding/json"
	"fmt"
	"glomers/spec/helpers"
	"glomers/src/broadcast"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func sendBroadcast(n helpers.Node, client string, message int) broadcast.BroadcastReply {
	msg := fmt.Sprintf(`{"type":"broadcast","message":%d}`, message)
	return n.Send[broadcast.BroadcastReply](client, json.RawMessage(msg))
}

func TestBroadcastHandler_RejectsMalformed(t *testing.T) {
	b := broadcast.NewBroadcast()
	node := helpers.Start(t, "n1", b)

	body := `{"type":"broadcast","message":"foo"}`
	resp := node.Send[broadcast.BroadcastReply]("c1", json.RawMessage(body))
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

	makeMsg := func(idx int) json.RawMessage {
		return json.RawMessage(fmt.Sprintf(`{"type":"broadcast","message":%d}`, idx))
	}
	node.SendConcurrently[broadcast.BroadcastReply](msgs, makeMsg)

	assert.Equal(t, msgs, len(b.Messages()))
}

func sendRead(n helpers.Node, client string) broadcast.ReadReply {
	return n.Send[broadcast.ReadReply](client, json.RawMessage(`{"type":"read"}`))
}

func TestReadHandlerReturnsNoMessages(t *testing.T) {
	b := broadcast.NewBroadcast()
	node := helpers.Start(t, "n1", b)
	resp := sendRead(node, "c1")
	assert.True(t, len(resp.Messages) == 0)
}

func TestReadHandlerReturnsMessages(t *testing.T) {
	b := broadcast.NewBroadcast()
	node := helpers.Start(t, "n1", b)
	sendBroadcast(node, "c1", 1)
	sendBroadcast(node, "c1", 2)
	resp := sendRead(node, "c1")
	assert.SliceContains(t, resp.Messages, 1)
	assert.SliceContains(t, resp.Messages, 2)
}

func TestReadHandler_PerformsConcurrentReads(t *testing.T) {
	msgs := 1000
	b := broadcast.NewBroadcast()
	node := helpers.Start(t, "n1", b)

	makeMsg := func(idx int) json.RawMessage {
		var body string
		if idx == 500 {
			body = `{"type":"read"}`
		} else {
			body = fmt.Sprintf(`{"type":"broadcast","message":%d}`, idx)
		}

		return json.RawMessage(body)
	}
	node.SendConcurrently[broadcast.BroadcastReply](msgs, makeMsg)

	assert.Equal(t, msgs-1, len(b.Messages()))
}

func TestTopologyHandler_RejectsMalformed(t *testing.T) {
	b := broadcast.NewBroadcast()
	node := helpers.Start(t, "n1", b)

	body := `{"type":"topology","topology":"foo"}`
	resp := node.Send[broadcast.TopologyReply]("c1", json.RawMessage(body))
	assert.Equal(t, "error", resp.Type)
}

func TestTopologyHandler_ReturnsOK(t *testing.T) {
	b := broadcast.NewBroadcast()
	node := helpers.Start(t, "n1", b)

	body := `{"type":"topology","topology":{}}`
	resp := node.Send[broadcast.TopologyReply]("c1", json.RawMessage(body))
	assert.Equal(t, "topology_ok", resp.Type)
}
