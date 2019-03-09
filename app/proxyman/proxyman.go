// Package proxyman defines applications for managing inbound and outbound proxies.
package proxyman

import (
	"context"

	"v2ray.com/core/common/session"
)

// ContextWithSniffingConfig is a wrapper of session.ContextWithContent.
// Deprecated. Use session.ContextWithContent directly.
func ContextWithSniffingConfig(ctx context.Context, c *SniffingConfig) context.Context {
	content := session.ContentFromContext(ctx)
	if content == nil {
		content = new(session.Content)
		ctx = session.ContextWithContent(ctx, content)
	}
	content.SniffingRequest.Enabled = c.Enabled
	content.SniffingRequest.OverrideDestinationForProtocol = c.DestinationOverride
	return ctx
}
