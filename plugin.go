package core

type PluginMetadata struct {
	Name string
}

const GetMetadataFuncName = "GetPluginMetadata"

type GetMetadataFunc func() PluginMetadata
