package conf

import (
	"encoding/json"

	"v2ray.com/core/common/errors"
)

type ConfigCreator func() interface{}

type ConfigCreatorCache map[string]ConfigCreator

func (v ConfigCreatorCache) RegisterCreator(id string, creator ConfigCreator) error {
	if _, found := v[id]; found {
		return errors.New("Config: ", id, " already registered.")
	}

	v[id] = creator
	return nil
}

func (v ConfigCreatorCache) CreateConfig(id string) (interface{}, error) {
	creator, found := v[id]
	if !found {
		return nil, errors.New("Config: Unknown config id: ", id)
	}
	return creator(), nil
}

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

func (v *JSONConfigLoader) LoadWithID(raw []byte, id string) (interface{}, error) {
	creator, found := v.cache[id]
	if !found {
		return nil, errors.New("Config: Unknown config id: ", id)
	}

	config := creator()
	if err := json.Unmarshal(raw, config); err != nil {
		return nil, err
	}
	return config, nil
}

func (v *JSONConfigLoader) Load(raw []byte) (interface{}, string, error) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil {
		return nil, "", err
	}
	rawID, found := obj[v.idKey]
	if !found {
		return nil, "", errors.New("Config: ", v.idKey, " not found in JSON context.")
	}
	var id string
	if err := json.Unmarshal(rawID, &id); err != nil {
		return nil, "", err
	}
	rawConfig := json.RawMessage(raw)
	if len(v.configKey) > 0 {
		configValue, found := obj[v.configKey]
		if !found {
			return nil, "", errors.New("Config: ", v.configKey, " not found in JSON content.")
		}
		rawConfig = configValue
	}
	config, err := v.LoadWithID([]byte(rawConfig), id)
	if err != nil {
		return nil, id, err
	}
	return config, id, nil
}
