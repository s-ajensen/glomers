package broadcast_test

import (
	"encoding/json"
	"fmt"
	"glomers/spec/helpers"
	"glomers/src/broadcast"
	"testing"

	"github.com/alecthomas/assert/v2"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func broadcastNode(t *testing.T, id string) (*broadcast.Broadcast, helpers.Node) {
	b := broadcast.NewBroadcast()
	return b, helpers.Start(t, id, nil, b)
}

func sendBroadcast(n helpers.Node, client string, message int) broadcast.BroadcastReply {
	msg := fmt.Sprintf(`{"type":"broadcast","message":%d}`, message)
	n.Send(client, json.RawMessage(msg))
	return n.ReadReply[broadcast.BroadcastReply]()
}

func TestBroadcastHandler_RejectsMalformed(t *testing.T) {
	_, node := broadcastNode(t, "n1")
	body := `{"type":"broadcast","message":"foo"}`
	node.Send("c1", json.RawMessage(body))
	resp := node.ReadReply[broadcast.BroadcastReply]()
	assert.Equal(t, "error", resp.Type)
}

func TestBroadcastHandlerReturnsOK(t *testing.T) {
	_, node := broadcastNode(t, "n1")
	assert.Equal(t, "broadcast_ok", sendBroadcast(node, "c1", 1).Type)
}

func TestBroadcastHandler_SchedulesConcurrentWrites(t *testing.T) {
	msgs := 1000
	b, node := broadcastNode(t, "n1")

	makeMsg := func(idx int) json.RawMessage {
		return json.RawMessage(fmt.Sprintf(`{"type":"broadcast","message":%d}`, idx))
	}
	node.SendConcurrently[broadcast.BroadcastReply](msgs, makeMsg)

	assert.Equal(t, msgs, len(b.Messages()))
}

func TestBroadcastHandler_DistributesMessage(t *testing.T) {
	b := broadcast.NewBroadcast()
	peers := []string{"n2"}
	node := helpers.Start(t, "n1", peers, b)
	node.Send("c1", json.RawMessage(`{"type":"broadcast","message":1}`))

	var dispatch maelstrom.Message
	var dispatchBody broadcast.BroadcastReq
	assert.NoError(t, node.Out.Decode(&dispatch))
	assert.NoError(t, json.Unmarshal(dispatch.Body, &dispatchBody))
	assert.Equal(t, 1, dispatchBody.Message)
}

func sendRead(n helpers.Node, client string) broadcast.ReadReply {
	n.Send(client, json.RawMessage(`{"type":"read"}`))
	return n.ReadReply[broadcast.ReadReply]()
}

func TestReadHandlerReturnsNoMessages(t *testing.T) {
	_, node := broadcastNode(t, "n1")
	resp := sendRead(node, "c1")
	assert.True(t, len(resp.Messages) == 0)
}

func TestReadHandlerReturnsMessages(t *testing.T) {
	_, node := broadcastNode(t, "n1")
	sendBroadcast(node, "c1", 1)
	sendBroadcast(node, "c1", 2)
	resp := sendRead(node, "c1")
	assert.SliceContains(t, resp.Messages, 1)
	assert.SliceContains(t, resp.Messages, 2)
}

func TestReadHandler_PerformsConcurrentReads(t *testing.T) {
	msgs := 1000
	b, node := broadcastNode(t, "n1")

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
	_, node := broadcastNode(t, "n1")

	body := `{"type":"topology","topology":"foo"}`
	node.Send("c1", json.RawMessage(body))
	resp := node.ReadReply[broadcast.TopologyReply]()
	assert.Equal(t, "error", resp.Type)
}

func TestTopologyHandler_ReturnsOK(t *testing.T) {
	_, node := broadcastNode(t, "n1")

	body := `{"type":"topology","topology":{}}`
	node.Send("c1", json.RawMessage(body))
	resp := node.ReadReply[broadcast.TopologyReply]()
	assert.Equal(t, "topology_ok", resp.Type)
}
