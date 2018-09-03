package configdb

import (
	"fmt"

	"github.com/hashicorp/go-memdb"
	"github.com/illarion/gonotify"
	log "github.com/sirupsen/logrus"
)

const (
	ConfigFolderPath     string = "/usr/local/zevenet/config"
	ConfigFilesIndexName string = "cfgFiles"
)

type ConfigDB struct {
	db          *memdb.MemDB
	schema      *memdb.DBSchema
	fileWatcher *gonotify.DirWatcher
}

var instance *ConfigDB

func GetInstance() *ConfigDB {
	if instance == nil {
		log.Warnf("Config DB not yet created, missing CreateConfigDb() call")
	}

	return instance
}

func CreateConfigDb() error {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{},
	}

	schema.Tables[TableNameGlobal] = getSchemaGlobal()

	db, err := memdb.NewMemDB(schema)
	if err != nil {
		return fmt.Errorf("Failed to initialize config DB: %v", err)
	}

	// start file watcher
	// (only writes and deletion of files affect the config DB)
	watcher, err := gonotify.NewDirWatcher(gonotify.IN_CLOSE_WRITE|gonotify.IN_DELETE, ConfigFolderPath)
	if err != nil {
		return fmt.Errorf("Failed to create file watcher: %v", err)
	}

	instance = &ConfigDB{
		db:          db,
		schema:      schema,
		fileWatcher: watcher,
	}

	go instance.watchFileChanges()

	return nil
}
