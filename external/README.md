# Note

The `external` dir in the project exists for historical reasons. We will migrate to use go mod to maintain 3rd party libraries.

Fow now, all modules under external are used by "quic-go" to support quic protocol, and can't migrate without a breaking change.

The plan is that we will remove the whole `external` dir when `quic-go` is tested to be matured enough in production.