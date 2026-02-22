// Package zerobot (host.go): installs protocol chain dispatch into ZeroBot. Call Install() from host (e.g. Lucy) after loading plugins.
package zerobot

import (
	"slices"
	"strconv"

	"github.com/Hafuunano/Protocol-ConvertTool/protocol"
	zero "github.com/wdvxdr1123/ZeroBot"
)

// makeContext builds protocol.Context from zero.Ctx with IsSuperAdmin, IsAdmin, and OnlyToMe wired.
// onlyToMe is true when handling zero.OnMessage(zero.OnlyToMe), false for zero.OnMessage().
func makeContext(ctx *zero.Ctx, onlyToMe bool) *Context {
	pc := NewContext(ctx)
	pc.IsSuperAdminFunc = func(userID string) bool {
		uid, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			return false
		}
		return slices.Contains(zero.BotConfig.SuperUsers, uid)
	}
	pc.IsAdminFunc = func() bool { return zero.AdminPermission(ctx) }
	pc.OnlyToMe = onlyToMe
	return pc
}

// Install registers ZeroBot message handlers with no global middlewares. Same as InstallWithMiddlewares(nil).
func Install() {
	InstallWithMiddlewares(nil)
}

// InstallWithMiddlewares registers ZeroBot message handlers, optionally wrapping the whole chain with middlewares.
// Middlewares are applied in order: mw[0](mw[1](...(dispatch(chain)))). Use for global whitelist etc., e.g. zerobot.InstallWithMiddlewares([]protocol.Middleware{whitelist.New(services.Cache)}).
func InstallWithMiddlewares(middlewares []protocol.Middleware) {
	wrap := func(h protocol.Handler) protocol.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			h = middlewares[i](h)
		}
		return h
	}
	// OnMessage: all messages -> default chain (HookMessage).
	zero.OnMessage().Handle(func(ctx *zero.Ctx) {
		wrap(protocol.Dispatch(protocol.Chain()))(makeContext(ctx, false))
	})
	// OnMessageReply: only reply to bot or @ bot -> HookMessageReply chain.
	zero.OnMessage(zero.OnlyToMe).Handle(func(ctx *zero.Ctx) {
		chain := protocol.ChainOn(protocol.HookMessageReply)
		if len(chain) > 0 {
			wrap(protocol.Dispatch(chain))(makeContext(ctx, true))
		}
	})
}
