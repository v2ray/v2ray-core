package core

// PluginMetadata contains some brief information regarding a plugin.
type PluginMetadata struct {
	// Name of the plugin
	Name string
}

// GetMetadataFuncName is the name of the function in the plugin to return PluginMetadata.
const GetMetadataFuncName = "GetPluginMetadata"

// GetMetadataFunc is the type of the function in the plugin to return PluginMetadata.
type GetMetadataFunc func() PluginMetadata

// LoadPlugins loads all possible plugins in the 'plugin' directory.
func LoadPlugins() error {
	return loadPluginsInternal()
}
