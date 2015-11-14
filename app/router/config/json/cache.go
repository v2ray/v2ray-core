package json

type ConfigObjectCreator func() interface{}

var (
	configCache map[string]ConfigObjectCreator
)

func RegisterRouterConfig(strategy string, creator ConfigObjectCreator) error {
	// TODO: check strategy
	configCache[strategy] = creator
	return nil
}

func CreateRouterConfig(strategy string) interface{} {
	creator, found := configCache[strategy]
	if !found {
		return nil
	}
	return creator()
}

func init() {
	configCache = make(map[string]ConfigObjectCreator)
}
