package session

import (
	"context"

	"v2ray.com/core/common/net"
	"v2ray.com/core/common/session"
	"v2ray.com/core/features/routing"
)

// Context is an implementation of routing.Context, which is a wrapper of context.context with session info.
type Context struct {
	Inbound  *session.Inbound
	Outbound *session.Outbound
	Content  *session.Content
}

// GetInboundTag implements routing.Context.
func (ctx *Context) GetInboundTag() string {
	if ctx.Inbound == nil {
		return ""
	}
	return ctx.Inbound.Tag
}

// GetSourceIPs implements routing.Context.
func (ctx *Context) GetSourceIPs() []net.IP {
	if ctx.Inbound == nil || !ctx.Inbound.Source.IsValid() {
		return nil
	}
	dest := ctx.Inbound.Source
	if dest.Address.Family().IsDomain() {
		return nil
	}

	return []net.IP{dest.Address.IP()}
}

// GetSourcePort implements routing.Context.
func (ctx *Context) GetSourcePort() net.Port {
	if ctx.Inbound == nil || !ctx.Inbound.Source.IsValid() {
		return 0
	}
	return ctx.Inbound.Source.Port
}

// GetTargetIPs implements routing.Context.
func (ctx *Context) GetTargetIPs() []net.IP {
	if ctx.Outbound == nil || !ctx.Outbound.Target.IsValid() {
		return nil
	}

	if ctx.Outbound.Target.Address.Family().IsIP() {
		return []net.IP{ctx.Outbound.Target.Address.IP()}
	}

	return nil
}

// GetTargetPort implements routing.Context.
func (ctx *Context) GetTargetPort() net.Port {
	if ctx.Outbound == nil || !ctx.Outbound.Target.IsValid() {
		return 0
	}
	return ctx.Outbound.Target.Port
}

// GetTargetDomain implements routing.Context.
func (ctx *Context) GetTargetDomain() string {
	if ctx.Outbound == nil || !ctx.Outbound.Target.IsValid() {
		return ""
	}
	dest := ctx.Outbound.Target
	if !dest.Address.Family().IsDomain() {
		return ""
	}
	return dest.Address.Domain()
}

// GetNetwork implements routing.Context.
func (ctx *Context) GetNetwork() net.Network {
	if ctx.Outbound == nil {
		return net.Network_Unknown
	}
	return ctx.Outbound.Target.Network
}

// GetProtocol implements routing.Context.
func (ctx *Context) GetProtocol() string {
	if ctx.Content == nil {
		return ""
	}
	return ctx.Content.Protocol
}

// GetUser implements routing.Context.
func (ctx *Context) GetUser() string {
	if ctx.Inbound == nil || ctx.Inbound.User == nil {
		return ""
	}
	return ctx.Inbound.User.Email
}

// GetAttributes implements routing.Context.
func (ctx *Context) GetAttributes() map[string]string {
	if ctx.Content == nil {
		return nil
	}
	return ctx.Content.Attributes
}

// AsRoutingContext creates a context from context.context with session info.
func AsRoutingContext(ctx context.Context) routing.Context {
	return &Context{
		Inbound:  session.InboundFromContext(ctx),
		Outbound: session.OutboundFromContext(ctx),
		Content:  session.ContentFromContext(ctx),
	}
}
