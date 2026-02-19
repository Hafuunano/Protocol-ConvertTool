package onebotv11

import (
	"encoding/json"
	"strconv"

	"github.com/Hafuunano/UniTransfer/protocol"
)

// Context implements protocol.Context for OneBot v11. Plugins use this via Plugin(ctx).
// Send/Reply do not perform network I/O; they build the OneBot v11 API request JSON and pass it to Out.
type Context struct {
	// Event is the current message event (set when handling an incoming message).
	Event *MessageEvent
	// Out receives the JSON payload that would be sent (action + params). Caller may log or send it elsewhere.
	// If nil, Send/Reply still build the payload but do nothing with it.
	Out func(payload []byte)
}

// Send implements protocol.Context. Builds send_private_msg/send_group_msg JSON and passes it to Out (no actual send).
func (c *Context) Send(msg protocol.Message) error {
	payload, err := c.buildSendPayload(msg, false)
	if err != nil {
		return err
	}
	if c.Out != nil && len(payload) > 0 {
		c.Out(payload)
	}
	return nil
}

// Reply implements protocol.Context. Builds send JSON with reply segment prepended and passes to Out (no actual send).
func (c *Context) Reply(msg protocol.Message) error {
	payload, err := c.buildSendPayload(msg, true)
	if err != nil {
		return err
	}
	if c.Out != nil && len(payload) > 0 {
		c.Out(payload)
	}
	return nil
}

// buildSendPayload returns the OneBot v11 API request JSON (action + params) for sending msg.
// If reply is true, prepends a reply segment using current event message_id.
func (c *Context) buildSendPayload(msg protocol.Message, reply bool) ([]byte, error) {
	if c == nil || c.Event == nil {
		return json.Marshal(map[string]any{"action": "send_private_msg", "params": map[string]any{"user_id": 0, "message": msg}})
	}
	m := msg
	if reply && c.Event.MessageID != 0 {
		replySeg := Reply(strconv.FormatInt(int64(c.Event.MessageID), 10))
		m = make(protocol.Message, 0, 1+len(msg))
		m = append(m, replySeg)
		m = append(m, msg...)
	}
	params := map[string]any{"message": []protocol.Segment(m)}
	if c.Event.MessageType == "group" && c.Event.GroupID != 0 {
		params["group_id"] = c.Event.GroupID
		return json.Marshal(map[string]any{"action": "send_group_msg", "params": params})
	}
	params["user_id"] = c.Event.UserID
	return json.Marshal(map[string]any{"action": "send_private_msg", "params": params})
}

// UserID implements protocol.Context (OneBot user_id).
func (c *Context) UserID() string {
	if c == nil || c.Event == nil {
		return ""
	}
	return strconv.FormatInt(c.Event.UserID, 10)
}

// GroupID implements protocol.Context (OneBot group_id, empty for private).
func (c *Context) GroupID() string {
	if c == nil || c.Event == nil || c.Event.GroupID == 0 {
		return ""
	}
	return strconv.FormatInt(c.Event.GroupID, 10)
}

// IncomingMessage implements protocol.Context (OneBot message).
func (c *Context) IncomingMessage() protocol.Message {
	if c == nil || c.Event == nil {
		return nil
	}
	return c.Event.Message
}

// PlainText implements protocol.Context (text-only, no CQ/media).
func (c *Context) PlainText() string {
	if c == nil || c.Event == nil {
		return ""
	}
	if len(c.Event.Message) > 0 {
		return ExtractPlainText(c.Event.Message)
	}
	return c.Event.RawMessage
}

// MessageID implements protocol.Context (OneBot message_id).
func (c *Context) MessageID() string {
	if c == nil || c.Event == nil {
		return ""
	}
	return strconv.FormatInt(int64(c.Event.MessageID), 10)
}

// RawMessage implements protocol.Context (OneBot raw_message).
func (c *Context) RawMessage() string {
	if c == nil || c.Event == nil {
		return ""
	}
	return c.Event.RawMessage
}

// SenderNickname implements protocol.Context (OneBot sender nickname/card).
func (c *Context) SenderNickname() string {
	if c == nil || c.Event == nil || c.Event.Sender == nil {
		return ""
	}
	return c.Event.Sender.DisplayName()
}

// Ensure Context implements protocol.Context at compile time.
var _ protocol.Context = (*Context)(nil)
