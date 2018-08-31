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

		if filename == "global.conf" {
			err := globalconfig.ReadZevenetGlobalConfig()
			if err != nil {
				log.Errorf("Failed to reload global config: %v", err)
			}
		} else {
			// clear related information from the config db
			t := db.db.Txn(true)
			defer t.Abort()

			t.DeleteAll("global", "cfgFiles", filename)

			t.Commit()
		}
	}
}
