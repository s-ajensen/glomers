package errors

type Error struct {
	Type string `json:"type"`
	Code int    `json:"code"`
	Text string `json:"text"`
}

func Malformed(msg string) Error {
	return Error{
		Type: "error",
		Code: 12,
		Text: msg,
	}
}
