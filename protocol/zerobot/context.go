package zerobot

import (
	"strconv"

	"github.com/Hafuunano/Protocol-ConvertTool/protocol"
	zero "github.com/wdvxdr1123/ZeroBot"
)

// Context wraps ZeroBot's *zero.Ctx and implements protocol.Context so plugins can use Plugin(ctx) with ZeroBot.
type Context struct {
	Ctx *zero.Ctx
	// IsSuperAdminFunc checks whether the given userID is a super admin. If nil, IsSuperAdmin() returns false.
	// Set this when building Context to enable super admin checks (e.g. from config SuperAdminIDs).
	IsSuperAdminFunc func(userID string) bool
	// IsAdminFunc checks whether the sender is group admin/owner. If nil, IsAdmin() returns false.
	// Host can set from zero.AdminPermission(ctx) or similar.
	IsAdminFunc func() bool
	// OnlyToMe is set by host: true when handling zero.OnMessage(zero.OnlyToMe), false for zero.OnMessage().
	OnlyToMe bool
	// blockNext is set by BlockNext(); host reads via ShouldBlockNext() to stop the chain.
	blockNext bool
}

// NewContext returns a protocol.Context that sends via ZeroBot's Ctx.
func NewContext(ctx *zero.Ctx) *Context {
	return &Context{Ctx: ctx}
}

// Send implements protocol.Context. It converts msg to ZeroBot message.Message and calls Ctx.Send.
func (c *Context) Send(msg protocol.Message) error {
	zbMsg := protocolMessageToZeroBot(msg)
	c.Ctx.Send(zbMsg)
	return nil
}

// Reply implements protocol.Context. ZeroBot Send in session context already replies.
func (c *Context) Reply(msg protocol.Message) error {
	return c.Send(msg)
}

// SendWithReply implements protocol.Context. Same as Reply for ZeroBot (session context replies).
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

// SendPoke implements protocol.Context. Sends a poke to the target user via NapCat send_poke API.
// Uses ZeroBot's Ctx.SendPoke(groupID, userID) which calls action "send_poke" with group_id and user_id (group_id 0 for private).
func (c *Context) SendPoke(targetUserID string) error {
	uid, err := strconv.ParseInt(targetUserID, 10, 64)
	if err != nil {
		return err
	}
	groupID := int64(0)
	if c.Ctx != nil && c.Ctx.Event != nil && c.Ctx.Event.GroupID != 0 {
		groupID = c.Ctx.Event.GroupID
	}
	c.Ctx.SendPoke(groupID, uid)
	return nil
}

// UserID implements protocol.Context.
func (c *Context) UserID() string {
	if c.Ctx == nil || c.Ctx.Event == nil {
		return ""
	}
	return strconv.FormatInt(c.Ctx.Event.UserID, 10)
}

// GroupID implements protocol.Context.
func (c *Context) GroupID() string {
	if c.Ctx == nil || c.Ctx.Event == nil || c.Ctx.Event.GroupID == 0 {
		return ""
	}
	return strconv.FormatInt(c.Ctx.Event.GroupID, 10)
}

// IncomingMessage implements protocol.Context.
func (c *Context) IncomingMessage() protocol.Message {
	if c.Ctx == nil || c.Ctx.Event == nil || c.Ctx.Event.Message == nil {
		return nil
	}
	return zeroBotMessageToProtocol(c.Ctx.Event.Message)
}

// PlainText implements protocol.Context.
func (c *Context) PlainText() string {
	if c.Ctx == nil {
		return ""
	}
	return c.Ctx.ExtractPlainText()
}

// MessageID implements protocol.Context (OneBot message_id).
func (c *Context) MessageID() string {
	if c.Ctx == nil || c.Ctx.Event == nil {
		return ""
	}
	return messageIDToString(c.Ctx.Event.RawMessageID)
}

// RawMessage implements protocol.Context (OneBot raw_message).
func (c *Context) RawMessage() string {
	if c.Ctx == nil || c.Ctx.Event == nil {
		return ""
	}
	return c.Ctx.Event.RawMessage
}

// SenderNickname implements protocol.Context (OneBot sender nickname/card).
func (c *Context) SenderNickname() string {
	if c.Ctx == nil || c.Ctx.Event == nil || c.Ctx.Event.Sender == nil {
		return ""
	}
	return c.Ctx.Event.Sender.Name()
}

// IsSuperAdmin implements protocol.Context. Returns true if current user is super admin (via IsSuperAdminFunc).
func (c *Context) IsSuperAdmin() bool {
	if c.IsSuperAdminFunc == nil {
		return false
	}
	return c.IsSuperAdminFunc(c.UserID())
}

// IsAdmin implements protocol.Context. Returns true if sender is group admin/owner (via IsAdminFunc, or false).
func (c *Context) IsAdmin() bool {
	if c.IsAdminFunc == nil {
		return false
	}
	return c.IsAdminFunc()
}

// IsOnlyToMe implements protocol.Context. Returns true when message is reply/@ bot; set by host (OnlyToMe field).
func (c *Context) IsOnlyToMe() bool {
	return c.OnlyToMe
}

// CommandPrefix implements protocol.Context. Returns ZeroBot CommandPrefix; empty defaults to "/".
func (c *Context) CommandPrefix() string {
	if zero.BotConfig.CommandPrefix != "" {
		return zero.BotConfig.CommandPrefix
	}
	return "/"
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
