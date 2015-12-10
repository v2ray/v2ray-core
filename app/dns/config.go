package dns

type CacheConfig interface {
	IsTrustedSource(tag string) bool
}
