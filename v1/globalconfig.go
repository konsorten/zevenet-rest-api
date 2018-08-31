package v1

import (
	"fmt"

	"gopkg.in/ini.v1"
)

const (
	GlobalConfigPath string = "/usr/local/zevenet/config/global.conf"
)

type ZevenetGlobalConfig struct {
	Version string
	ZAPIKey string
}

func ReadZevenetGlobalConfig() (*ZevenetGlobalConfig, error) {
	cfg, err := ini.Load(GlobalConfigPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to load global config: %v", GlobalConfigPath)
	}

	return &ZevenetGlobalConfig{
		Version: cfg.Section("").Key("$version").String(),
		ZAPIKey: cfg.Section("").Key("$zapikey").String(),
	}, nil
}
