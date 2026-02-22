// Package protocol (hooks.go): hook name constants for multi-chain registration.
// Host and plugins use these to register or obtain chains by event type.
package protocol

// Hook names for post_type-level events (OneBot v11: message, notice, request, meta_event).
const (
	HookMessage   = "message"    // all messages; default for Register() and Chain()
	HookNotice    = "notice"
	HookRequest   = "request"
	HookMetaEvent = "meta_event"
)

// Message sub-type hooks: finer-grained message chains (e.g. OnMessage vs OnMessageReply).
const (
	// HookMessageReply: only messages that reply to bot or @ bot (e.g. zero.OnMessage(zero.OnlyToMe)).
	HookMessageReply = "message_reply"
	// Optional extensions (add when needed): HookMessagePrivate = "message_private", HookMessageGroup = "message_group"
)
