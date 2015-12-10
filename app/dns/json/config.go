package json

import (
	"strings"

	serialjson "github.com/v2ray/v2ray-core/common/serial/json"
)

type TagList map[string]bool

func NewTagList(tags []string) TagList {
	list := TagList(make(map[string]bool))
	for _, tag := range tags {
		list[strings.TrimSpace(tag)] = true
	}
	return list
}

func (this *TagList) UnmarshalJSON(data []byte) error {
	tags, err := serialjson.UnmarshalStringList(data)
	if err != nil {
		return err
	}
	*this = NewTagList(tags)
	return nil
}

type CacheConfig struct {
	TrustedTags TagList `json:"trustedTags"`
}

func (this *CacheConfig) IsTrustedSource(tag string) bool {
	_, found := this.TrustedTags[tag]
	return found
}
