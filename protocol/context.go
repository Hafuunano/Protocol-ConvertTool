// Package protocol (context.go): plugin context interface.
// Plugins use Context for all sending and protocol operations; the implementation is provided by each protocol (e.g. onebotv11, zerobot).
package protocol

type Context interface {
	// Send sends a unified message. The implementation translates it to the protocol format and performs the send.
	Send(msg Message) error
	// Reply sends a reply to the current message (e.g. prepends reply segment or uses reply API).
	Reply(msg Message) error
	// UserID returns the sender user ID of the current event.
	UserID() string
	// GroupID returns the group ID when in a group chat, or empty/"0" for private chat.
	GroupID() string
	// IncomingMessage returns the current event message as unified Message.
	IncomingMessage() Message
	// PlainText returns the plain text of the current message (no CQ/media).
	PlainText() string
	// MessageID returns the current event message_id (for reply/recall).
	MessageID() string
	// RawMessage returns the original raw_message string from the event.
	RawMessage() string
	// SenderNickname returns the sender display name (nickname or card per OneBot sender).
	SenderNickname() string
}
