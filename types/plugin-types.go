package types

import "github.com/gin-gonic/gin"

type PluginEnableConfig struct {
	Enable  bool               `yaml:"enable"`
	Plugins []PluginBaseConfig `yaml:"plugins"`
}

type PluginBaseConfig struct {
	Name           string `yaml:"name"`
	Token          string `yaml:"token"`
	Route          string `yaml:"route"`
	PluginPath     string `yaml:"plugin_path"`
	ConfigPath     string `yaml:"config_path"`
	ShareWhiteList bool   `yaml:"share_white_list"`
}

type InitPlugin func(string) (func(c *gin.Context), error)
