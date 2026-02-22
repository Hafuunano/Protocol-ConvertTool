// Package protocol (handler.go): handler and middleware types for plugin chain.
// Used by host (e.g. Lucy) to build the middleware chain and invoke plugins.
package protocol

// Handler represents a single "handle one message" unit. Plugins and the final
// dispatcher are invoked as Handler. Host calls Handler(ctx) for each message.
type Handler = func(Context)

// Middleware wraps the next Handler and returns a new Handler. The returned
// Handler should run middleware logic then call next(ctx) to continue the chain.
// Example: func Whitelist(next Handler) Handler { return func(ctx Context) { if ok { next(ctx) } } }
type Middleware = func(next Handler) Handler
