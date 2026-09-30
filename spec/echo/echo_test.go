package echo_test

import (
	"encoding/json"
	"glomers/spec/helpers"
	"glomers/src/echo"
	"testing"

	"github.com/alecthomas/assert/v2"
)

type echoReply struct {
	Type string `json:"type"`
	Echo string `json:"echo"`
}

func TestHandlerReturns_Message(t *testing.T) {
	e := echo.Echo{}
	node1 := helpers.Start(t, "n1", e)

	resp := node1.Send[echoReply]("c1", json.RawMessage(`{"type":"echo","echo":"foo"}`))
	assert.Equal(t, "echo_ok", resp.Type)
	assert.Equal(t, "foo", resp.Echo)
}
