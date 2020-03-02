package sonnarplugins

import (
	"fmt"

	"github.com/diegocsandrim/engagement-by-quality/cmd"
)

type SonnarPlugin interface {
	Setup() error
	Run() error
	ReportFlag() string
}

func CreatePlugins(cmdFactory *cmd.CmdFactory, pluginNames []string) ([]SonnarPlugin, error) {
	plugins := make([]SonnarPlugin, 0, len(pluginNames))

	for _, pluginName := range pluginNames {
		plugin, err := createPlugin(cmdFactory, pluginName)
		if err != nil {
			return nil, err
		}

		plugins = append(plugins, plugin)
	}

	for _, plugin := range plugins {
		err := plugin.Setup()
		if err != nil {
			return nil, err
		}
	}

	return plugins, nil
}

func createPlugin(cmdFactory *cmd.CmdFactory, pluginName string) (SonnarPlugin, error) {
	switch pluginName {
	case "golangci-lint":
		plugin := GolangCILintPlugin{
			cmdFactory: cmdFactory,
		}
		return &plugin, nil
	default:
		return nil, fmt.Errorf("unsupported plugin %s", pluginName)
	}
}
