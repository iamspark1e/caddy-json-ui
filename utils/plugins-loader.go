package utils

// https://github.com/bigwhite/experiments/blob/master/go-plugin/demo1/pkg/pkg1/pkg1.go

import (
	"log"
	"plugin"
)

func LoadAndInvokeSomethingFromPlugin(pluginPath string) (plugin.Symbol, error) {
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, err
	}

	// 导出函数变量
	f, err := p.Lookup("InitPlugin")
	if err != nil {
		log.Printf("Lookup error: %s", err.Error())
		return nil, err
	}

	return f, nil
}
