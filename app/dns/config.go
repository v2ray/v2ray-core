package dns

type CacheConfig interface {
	TrustedSource() []string
}
