package zerobot

import (
	"strconv"

	"github.com/Hafuunano/UniTransfer/protocol"
	zero "github.com/wdvxdr1123/ZeroBot"
)

// Context wraps ZeroBot's *zero.Ctx and implements protocol.Context so plugins can use Plugin(ctx) with ZeroBot.
type Context struct {
	Ctx *zero.Ctx
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

// Ensure Context implements protocol.Context at compile time.
var _ protocol.Context = (*Context)(nil)
