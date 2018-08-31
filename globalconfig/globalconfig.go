package globalconfig

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	flock "github.com/theckman/go-flock"
	"gopkg.in/ini.v1"
)

const (
	GlobalConfigPath string = "/usr/local/zevenet/config/global.conf"
)

type ZevenetGlobalConfig struct {
	Version string
	ZAPIKey string
}

var instance *ZevenetGlobalConfig

func GetZevenetGlobalConfig() *ZevenetGlobalConfig {
	if instance == nil {
		log.Warnf("Zevenet global config not yet created, missing ReadZevenetGlobalConfig() call")
	}

	return instance
}

func ReadZevenetGlobalConfig() error {
	log.Debugf("Loading global configuration file: %v", GlobalConfigPath)

	// lock the file
	fileLock := flock.NewFlock(GlobalConfigPath)
	locked, err := fileLock.TryRLock()
	if err != nil {
		return fmt.Errorf("Failed to lock global config: %v: %v", GlobalConfigPath, err)
	}
	if locked {
		defer fileLock.Unlock()
	}

	// read file
	cfg, err := ini.Load(GlobalConfigPath)
	if err != nil {
		return fmt.Errorf("Failed to load global config: %v: %v", GlobalConfigPath, err)
	}

	instance = &ZevenetGlobalConfig{
		Version: cfg.Section("").Key("$version").String(),
		ZAPIKey: cfg.Section("").Key("$zapikey").String(),
	}

	return nil
}
