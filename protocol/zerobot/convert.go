// Package zerobot provides protocol.Context implementation by wrapping ZeroBot's *zero.Ctx.
package zerobot

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Hafuunano/Protocol-ConvertTool/protocol"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// protocolMessageToZeroBot converts unified protocol.Message to ZeroBot message.Message.
func protocolMessageToZeroBot(msg protocol.Message) message.Message {
	if len(msg) == 0 {
		return nil
	}
	out := make(message.Message, 0, len(msg))
	for _, seg := range msg {
		out = append(out, protocolSegmentToZeroBot(seg))
	}
	return out
}

func protocolSegmentToZeroBot(seg protocol.Segment) message.Segment {
	data := make(map[string]string)
	for k, v := range seg.Data {
		data[k] = dataToString(v)
	}
	return message.Segment{
		Type: seg.Type,
		Data: data,
	}
}

// zeroBotMessageToProtocol converts ZeroBot message.Message to unified protocol.Message.
func zeroBotMessageToProtocol(zbMsg message.Message) protocol.Message {
	if len(zbMsg) == 0 {
		return nil
	}
	out := make(protocol.Message, 0, len(zbMsg))
	for _, seg := range zbMsg {
		out = append(out, zeroBotSegmentToProtocol(seg))
	}
	return out
}

func zeroBotSegmentToProtocol(seg message.Segment) protocol.Segment {
	data := make(map[string]interface{})
	for k, v := range seg.Data {
		data[k] = v
	}
	return protocol.Segment{
		Type: seg.Type,
		Data: data,
	}
}

// messageIDToString converts OneBot message_id (JSON number or string) to string.
func messageIDToString(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return ""
	}
	switch x := v.(type) {
	case float64:
		return strconv.FormatInt(int64(x), 10)
	case string:
		return x
	default:
		return fmt.Sprint(v)
	}
}

// dataToString converts segment data value to string (same logic as onebotv11 segmentDataToString).
func dataToString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case int:
		return strconv.Itoa(x)
	case int64:
		return strconv.FormatInt(x, 10)
	case bool:
		return strconv.FormatBool(x)
	default:
		return fmt.Sprint(x)
	}
}
