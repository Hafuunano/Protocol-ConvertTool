// Package protocol (registry.go): global plugin chain registration.
// Plugins call Register(Plugin) in init(); host uses Chain() to get the list for dispatch.
package protocol

// defaultChain holds handlers registered by plugins via Register.
var defaultChain []Handler

// Register appends a handler to the default chain. Called by plugins in init().
func Register(h Handler) {
	defaultChain = append(defaultChain, h)
}

// Chain returns a copy of the registered handler list for dispatch.
func Chain() []Handler {
	if len(defaultChain) == 0 {
		return nil
	}
	out := make([]Handler, len(defaultChain))
	copy(out, defaultChain)
	return out
}
