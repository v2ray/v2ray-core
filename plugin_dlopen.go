// +build linux,cgo darwin,cgo

package core

import (
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"v2ray.com/core/common/platform"
)

func loadPluginsInternal() error {
	pluginPath := platform.GetPluginDirectory()

	dir, err := os.Open(pluginPath)
	if err != nil {
		return err
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".so") {
			p, err := plugin.Open(filepath.Join(pluginPath, file.Name()))
			if err != nil {
				return err
			}
			f, err := p.Lookup(GetMetadataFuncName)
			if err != nil {
				return err
			}
			if gmf, ok := f.(func() PluginMetadata); ok {
				metadata := GetMetadataFunc(gmf)()
				newError("plugin (", metadata.Name, ") loaded.").WriteToLog()
			} else {
				return newError(file.Name(), " is not a valid plugin.")
			}
		}
	}

	return nil
}
