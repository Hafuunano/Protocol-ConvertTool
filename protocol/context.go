// Package protocol (context.go): plugin context interface.
// Plugins use Context for all sending and protocol operations; the implementation is provided by each protocol (e.g. onebotv11, zerobot).
package protocol

type Context interface {
	// Send sends a unified message. The implementation translates it to the protocol format and performs the send.
	Send(msg Message) error
	// Reply sends a reply to the current message (e.g. prepends reply segment or uses reply API).
	Reply(msg Message) error
	// SendWithReply sends the given message as a reply to the current message (same semantics as Reply).
	SendWithReply(msg Message) error
	// SendPlainMessage sends a single plain-text message. Convenience for sending text without building Message.
	SendPlainMessage(text string) error
	// SendWithImage sends a single image. file can be path, URL, or base64://... per protocol.
	SendWithImage(file string) error
	// SendWithImageAndText sends one image followed by plain text (e.g. image + caption). file can be path, URL, or base64://...
	SendWithImageAndText(file string, text string) error
	// SendPoke sends a poke to the target user (NapCat/OneBot extension: send_poke). Group context includes group_id; private uses user_id only.
	SendPoke(targetUserID string) error
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
	// IsSuperAdmin returns true if the current user is a super admin (e.g. configured bot owner/admin).
	// Plugins can use this to gate privileged commands (TriggerOnlySuperAdmin).
	IsSuperAdmin() bool
	// IsAdmin returns true if the sender is a group admin or owner (TriggerOnlyAdmin). Private chat has no role, returns false.
	IsAdmin() bool
	// IsOnlyToMe returns true when the message is reply-to-bot or @-bot (OnlyToMe). Set by host or derived from message.
	IsOnlyToMe() bool
	// CommandPrefix returns the bot command prefix (e.g. "/", "!"). Use with command names for OnCommand-style matching.
	CommandPrefix() string
	// BlockNext marks that this plugin handled the event; host should not run the rest of the chain (BlockNextPlugin).
	BlockNext()
	// ShouldBlockNext returns true after BlockNext() was called. Used by host after each handler to decide whether to break.
	ShouldBlockNext() bool
}
