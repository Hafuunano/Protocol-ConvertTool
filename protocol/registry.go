// Package protocol (registry.go): plugin chain registration by hook.
// Plugins call Register(Plugin) or RegisterOn(hook, Plugin) in init(); host uses Chain() or ChainOn(hook) for dispatch.
package protocol

// chains holds handlers per hook. Key = hook name (e.g. HookMessage, HookMessageReply), value = handler list.
var chains = make(map[string][]Handler)

// Register appends a handler to the default chain (HookMessage). Called by plugins in init(). Backward compatible.
func Register(h Handler) {
	RegisterOn(HookMessage, h)
}

// RegisterOn appends a handler to the chain for the given hook. Use for message_reply, notice, request, etc.
func RegisterOn(hook string, h Handler) {
	chains[hook] = append(chains[hook], h)
}

// Chain returns a copy of the default (HookMessage) chain for dispatch. Backward compatible.
func Chain() []Handler {
	return ChainOn(HookMessage)
}

// ChainOn returns a copy of the registered handler list for the given hook. Returns nil if none registered.
func ChainOn(hook string) []Handler {
	list := chains[hook]
	if len(list) == 0 {
		return nil
	}
	out := make([]Handler, len(list))
	copy(out, list)
	return out
}
