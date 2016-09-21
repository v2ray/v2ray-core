// +build json

package loader

import (
	"encoding/json"

	"v2ray.com/core/common"
	"v2ray.com/core/common/log"
)

type JSONConfigLoader struct {
	cache     ConfigCreatorCache
	idKey     string
	configKey string
}

func NewJSONConfigLoader(cache ConfigCreatorCache, idKey string, configKey string) *JSONConfigLoader {
	return &JSONConfigLoader{
		idKey:     idKey,
		configKey: configKey,
		cache:     cache,
	}
}

func (this *JSONConfigLoader) LoadWithID(raw []byte, id string) (interface{}, error) {
	config, err := this.cache.CreateConfig(id)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(raw, config); err != nil {
		return nil, err
	}
	return config, nil
}

func (this *JSONConfigLoader) Load(raw []byte) (interface{}, string, error) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil {
		return nil, "", err
	}
	rawID, found := obj[this.idKey]
	if !found {
		log.Error(this.idKey, " not found in JSON content.")
		return nil, "", common.ErrObjectNotFound
	}
	var id string
	if err := json.Unmarshal(rawID, &id); err != nil {
		return nil, "", err
	}
	rawConfig := json.RawMessage(raw)
	if len(this.configKey) > 0 {
		configValue, found := obj[this.configKey]
		if !found {
			log.Error(this.configKey, " not found in JSON content.")
			return nil, "", common.ErrObjectNotFound
		}
		rawConfig = configValue
	}
	config, err := this.LoadWithID([]byte(rawConfig), id)
	if err != nil {
		return nil, id, err
	}
	return config, id, nil
}
