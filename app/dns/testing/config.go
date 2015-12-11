package testing

type CacheConfig struct {
	TrustedTags map[string]bool
}

func (this *CacheConfig) IsTrustedSource(tag string) bool {
	_, found := this.TrustedTags[tag]
	return found
}
