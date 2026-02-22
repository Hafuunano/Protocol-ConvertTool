package onebotv11

import (
	"encoding/json"
	"strconv"

	"github.com/Hafuunano/Protocol-ConvertTool/protocol"
)

// Context implements protocol.Context for OneBot v11. Plugins use this via Plugin(ctx).
// Send/Reply do not perform network I/O; they build the OneBot v11 API request JSON and pass it to Out.
type Context struct {
	// Event is the current message event (set when handling an incoming message).
	Event *MessageEvent
	// Out receives the JSON payload that would be sent (action + params). Caller may log or send it elsewhere.
	// If nil, Send/Reply still build the payload but do nothing with it.
	Out func(payload []byte)
	// IsSuperAdminFunc checks whether the given userID is a super admin. If nil, IsSuperAdmin() returns false.
	// Set this when building Context to enable super admin checks (e.g. from config SuperAdminIDs).
	IsSuperAdminFunc func(userID string) bool
	// OnlyToMe, if non-nil, overrides IsOnlyToMe() result. Otherwise derived from message (at-segment to self).
	OnlyToMe *bool
	// commandPrefix is the bot command prefix for OnCommand-style matching. If empty, CommandPrefix() returns "/".
	commandPrefix string
	// blockNext is set by BlockNext(); host reads via ShouldBlockNext() to stop the chain.
	blockNext bool
}

// Send implements protocol.Context. Builds send_private_msg/send_group_msg JSON and passes it to Out (no actual send).
// If msg is a single poke segment, uses SendPoke (send_poke API) for NapCat compatibility.
func (c *Context) Send(msg protocol.Message) error {
	if len(msg) == 1 && msg[0].Type == protocol.SegmentTypePoke {
		if id := pokeTargetFromSegment(msg[0]); id != "" {
			return c.SendPoke(id)
		}
	}
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

// SendWithReply implements protocol.Context. Same as Reply: sends msg with reply segment prepended.
func (c *Context) SendWithReply(msg protocol.Message) error {
	return c.Reply(msg)
}

// SendPlainMessage implements protocol.Context. Sends a single plain-text message.
func (c *Context) SendPlainMessage(text string) error {
	return c.Send(protocol.Message{
		{Type: protocol.SegmentTypeText, Data: map[string]any{"text": text}},
	})
}

// SendWithImage implements protocol.Context. Sends a single image; file can be path, URL, or base64://...
func (c *Context) SendWithImage(file string) error {
	return c.Send(protocol.Message{
		{Type: protocol.SegmentTypeImage, Data: map[string]any{"file": file}},
	})
}

// SendWithImageAndText implements protocol.Context. Sends image then text (e.g. image + caption).
func (c *Context) SendWithImageAndText(file string, text string) error {
	return c.Send(protocol.Message{
		{Type: protocol.SegmentTypeImage, Data: map[string]any{"file": file}},
		{Type: protocol.SegmentTypeText, Data: map[string]any{"text": text}},
	})
}

// SendPoke implements protocol.Context. Uses NapCat/OneBot extension send_poke API (group: user_id + group_id; private: user_id).
func (c *Context) SendPoke(targetUserID string) error {
	uid, err := strconv.ParseInt(targetUserID, 10, 64)
	if err != nil {
		return err
	}
	params := map[string]any{"user_id": uid}
	if c.Event != nil && c.Event.MessageType == "group" && c.Event.GroupID != 0 {
		params["group_id"] = c.Event.GroupID
	}
	payload, err := json.Marshal(map[string]any{"action": "send_poke", "params": params})
	if err != nil {
		return err
	}
	if c.Out != nil && len(payload) > 0 {
		c.Out(payload)
	}
	return nil
}

// pokeTargetFromSegment returns target user ID from a poke segment Data (id or qq).
func pokeTargetFromSegment(seg protocol.Segment) string {
	if seg.Data == nil {
		return ""
	}
	for _, k := range []string{"qq", "id"} {
		if v, ok := seg.Data[k]; ok && v != nil {
			switch x := v.(type) {
			case string:
				if x != "" {
					return x
				}
			case float64:
				return strconv.FormatInt(int64(x), 10)
			case int:
				return strconv.Itoa(x)
			case int64:
				return strconv.FormatInt(x, 10)
			}
		}
	}
	return ""
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

// IsSuperAdmin implements protocol.Context. Returns true if current user is super admin (via IsSuperAdminFunc).
func (c *Context) IsSuperAdmin() bool {
	if c == nil || c.IsSuperAdminFunc == nil {
		return false
	}
	return c.IsSuperAdminFunc(c.UserID())
}

// IsAdmin implements protocol.Context. Returns true if sender is group admin or owner (Sender.Role).
func (c *Context) IsAdmin() bool {
	if c == nil || c.Event == nil || c.Event.Sender == nil {
		return false
	}
	r := c.Event.Sender.Role
	return r == "admin" || r == "owner"
}

// CommandPrefix implements protocol.Context. Returns configured prefix or "/" if not set.
func (c *Context) CommandPrefix() string {
	if c != nil && c.commandPrefix != "" {
		return c.commandPrefix
	}
	return "/"
}

// IsOnlyToMe implements protocol.Context. If OnlyToMe is set by host, use it; else true when message @-s self.
func (c *Context) IsOnlyToMe() bool {
	if c == nil || c.Event == nil {
		return false
	}
	if c.OnlyToMe != nil {
		return *c.OnlyToMe
	}
	selfIDStr := strconv.FormatInt(c.Event.SelfID, 10)
	for _, seg := range c.Event.Message {
		if seg.Type == protocol.SegmentTypeAt {
			if qq, _ := seg.Data["qq"]; qq != nil {
				switch v := qq.(type) {
				case string:
					if v == selfIDStr {
						return true
					}
				case float64:
					if int64(v) == c.Event.SelfID {
						return true
					}
				}
			}
		}
	}
	return false
}

// BlockNext implements protocol.Context. Marks that the plugin handled the event; host should stop the chain.
func (c *Context) BlockNext() {
	c.blockNext = true
}

// ShouldBlockNext implements protocol.Context. Returns true after BlockNext() was called.
func (c *Context) ShouldBlockNext() bool {
	return c.blockNext
}

// Ensure Context implements protocol.Context at compile time.
var _ protocol.Context = (*Context)(nil)
