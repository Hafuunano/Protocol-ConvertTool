// Package protocol (engine.go): fluent registration API.
// Meta is required at init: use Engine.WithMeta(Meta) then chain OnMessage()/OnMessageReply(), IsOnlyToMe(), etc., then Func(handler).
package protocol

import "strings"

// Engine is the entry for fluent handler registration. Plugins must call Engine.WithMeta(Meta) first, then OnMessage().Func(handler), etc.
var Engine = &engine{}

type engine struct{}

// PluginBuilder is returned by Engine.WithMeta(meta). Use it to chain OnMessage()/OnMessageReply() then conditions and Func(handler).
type PluginBuilder struct {
	meta any
}

// WithMeta is required at plugin init. Pass plugin metadata (e.g. types.PluginEngine); may be nil for plugins without meta. Returns PluginBuilder for chain calls.
func (e *engine) WithMeta(meta interface{}) *PluginBuilder {
	return &PluginBuilder{meta: meta}
}

// OnMessage returns a Builder for the default message chain (HookMessage).
// With no args: all messages trigger. With one arg: only when PlainText (trimmed) equals that string, e.g. OnMessage("ping") replaces ctx.PlainText() == "ping" in handler.
func (b *PluginBuilder) OnMessage(trigger ...string) *Builder {
	exact := ""
	if len(trigger) >= 1 {
		exact = trigger[0]
	}
	return &Builder{hook: HookMessage, exactText: exact}
}

// OnMessageNamed returns a Builder for a named message chain. Name "" or "Global" maps to HookMessage; other names reserved for future use.
func (b *PluginBuilder) OnMessageNamed(name string) *Builder {
	if name == "" || name == "Global" {
		return &Builder{hook: HookMessage}
	}
	return &Builder{hook: HookMessage}
}

// OnMessageReply returns a Builder for the reply/@-bot-only chain (HookMessageReply).
func (b *PluginBuilder) OnMessageReply() *Builder {
	return &Builder{hook: HookMessageReply}
}

// Builder is returned by PluginBuilder.OnMessage / OnMessageReply. Configure trigger conditions then call Func(handler).
type Builder struct {
	hook           string
	exactText      string // when non-empty, handler runs only when PlainText (trimmed) equals exactText
	onlyToMe       bool
	onlyAdmin      bool
	onlySuperAdmin bool
}

// IsOnlyToMe restricts the handler to run only when the message is reply-to-bot or @-bot.
func (b *Builder) IsOnlyToMe() *Builder {
	b.onlyToMe = true
	return b
}

// IsOnlyAdmin restricts the handler to run only when the sender is group admin or owner.
func (b *Builder) IsOnlyAdmin() *Builder {
	b.onlyAdmin = true
	return b
}

// IsOnlySuperAdmin restricts the handler to run only when the sender is super admin.
func (b *Builder) IsOnlySuperAdmin() *Builder {
	b.onlySuperAdmin = true
	return b
}

// Func registers the handler on the builder's hook. If exactText is set, OnlyToMe, OnlyAdmin, or OnlySuperAdmin are set,
// the handler is wrapped so it is only invoked when all those conditions pass.
func (b *Builder) Func(h Handler) {
	wrapped := func(ctx Context) {
		if b.exactText != "" && strings.TrimSpace(ctx.PlainText()) != b.exactText {
			return
		}
		if b.onlyToMe && !ctx.IsOnlyToMe() {
			return
		}
		if b.onlyAdmin && !ctx.IsAdmin() {
			return
		}
		if b.onlySuperAdmin && !ctx.IsSuperAdmin() {
			return
		}
		h(ctx)
	}
	RegisterOn(b.hook, wrapped)
}
