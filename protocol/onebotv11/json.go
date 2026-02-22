package onebotv11

import (
	"encoding/json"

	"github.com/Hafuunano/Protocol-ConvertTool/protocol"
)

// MarshalMessage encodes msg to OneBot v11 array format JSON: [{"type":"...","data":{...}}, ...].
func MarshalMessage(msg protocol.Message) ([]byte, error) {
	return json.Marshal([]protocol.Segment(msg))
}

// UnmarshalMessage decodes OneBot v11 array format JSON into a Message.
func UnmarshalMessage(data []byte) (protocol.Message, error) {
	var segs []protocol.Segment
	if err := json.Unmarshal(data, &segs); err != nil {
		return nil, err
	}
	return protocol.Message(segs), nil
}

// UnmarshalMessageFromString is like UnmarshalMessage but accepts a string.
func UnmarshalMessageFromString(s string) (protocol.Message, error) {
	return UnmarshalMessage([]byte(s))
}
