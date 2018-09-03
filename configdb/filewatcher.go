package configdb

import (
	"path"

	"github.com/konsorten/zevenet-rest-api/globalconfig"
	log "github.com/sirupsen/logrus"
)

func (db *ConfigDB) watchFileChanges() {
	for {
		ev := <-db.fileWatcher.C
		if ev.Eof {
			break
		}

		filename := path.Base(ev.Name)

		log.Infof("Configuration file changed: %v", filename)

		// special handling for global.conf
		if filename == "global.conf" {
			err := globalconfig.ReadZevenetGlobalConfig()
			if err != nil {
				log.Errorf("Failed to reload global config: %v", err)
			}
		}

		// clear related information from the config db
		t := db.db.Txn(true)
		defer t.Abort()
		var deleted int

		for _, table := range db.schema.Tables {
			if _, ok := table.Indexes[ConfigFilesIndexName]; ok {
				d, err := t.DeleteAll(table.Name, ConfigFilesIndexName, filename)
				if err != nil {
					log.Warnf("Failed to delete from '%v' table with index '%v' from config DB", table.Name, ConfigFilesIndexName)
					continue
				}
				deleted += d
			}
		}

		log.Debugf("Entries purged from cache/config DB: %v", deleted)

		t.Commit()
	}
}
