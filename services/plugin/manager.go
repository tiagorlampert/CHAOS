package services

import (
    "CHAOS/entities"
)

type PluginManager struct {
    Plugins []entities.Plugin
    Loader  *PluginLoader
}

func NewPluginManager() *PluginManager {
    return &PluginManager{
        Loader: &PluginLoader{},
    }
}

func (pm *PluginManager) LoadPlugins(pluginDir string) error {
    plugins, err := pm.Loader.LoadPlugins(pluginDir)
    if err != nil {
        return err
    }
    pm.Plugins = plugins
    return nil
}

func (pm *PluginManager) ExecutePlugins() error {
    for _, p := range pm.Plugins {
        if err := p.Execute(); err != nil {
            return err
        }
    }
    return nil
}

func (pm *PluginManager) CleanupPlugins() error {
    for _, p := range pm.Plugins {
        if err := p.Cleanup(); err != nil {
            return err
        }
    }
    return nil
}