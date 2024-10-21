package services

import (
    "plugin"
    "path/filepath"
    "CHAOS/entities"
)

type PluginLoader struct {}

func (pl *PluginLoader) LoadPlugins(pluginDir string) ([]entities.Plugin, error) {
    var plugins []entities.Plugin

    files, err := filepath.Glob(filepath.Join(pluginDir, "*.so"))
    if err != nil {
        return nil, err
    }

    for _, file := range files {
        p, err := plugin.Open(file)
        if err != nil {
            return nil, err
        }

        symPlugin, err := p.Lookup("Plugin")
        if err != nil {
            return nil, err
        }

        var plug entities.Plugin
        plug, ok := symPlugin.(entities.Plugin)
        if !ok {
            return nil, fmt.Errorf("unexpected type from module symbol")
        }

        if err := plug.Initialize(); err != nil {
            return nil, err
        }

        plugins = append(plugins, plug)
    }

    return plugins, nil
}