package configdb

import (
	"github.com/Jeffail/gabs"
	"github.com/hashicorp/go-memdb"
	log "github.com/sirupsen/logrus"
)

const (
	TableNameGlobal string = "global"
)

type GlobalObject struct {
	Name        string
	ConfigFiles []string
	JSON        []byte
}

func (obj *GlobalObject) Object() *gabs.Container {
	r, err := gabs.ParseJSON(obj.JSON)
	if err != nil {
		log.Errorf("Failed to re-parse JSON from cache: %v", obj.Name)
		return gabs.New()
	}
	return r
}

func getSchemaGlobal() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: TableNameGlobal,
		Indexes: map[string]*memdb.IndexSchema{
			"id": &memdb.IndexSchema{
				Name:    "id",
				Unique:  true,
				Indexer: &memdb.StringFieldIndex{Field: "Name"},
			},
			ConfigFilesIndexName: &memdb.IndexSchema{
				Name:         ConfigFilesIndexName,
				Indexer:      &memdb.StringSliceFieldIndex{Field: "ConfigFiles"},
				AllowMissing: true,
			},
		},
	}
}

func (db *ConfigDB) AddGlobal(name string, obj *gabs.Container, cfgFiles ...string) error {
	log.Debugf("Cache ADD on global: %v", name)

	t := db.db.Txn(true)
	defer t.Abort()

	err := t.Insert(TableNameGlobal, &GlobalObject{
		Name:        name,
		JSON:        obj.Bytes(),
		ConfigFiles: cfgFiles,
	})
	if err != nil {
		return err
	}

	t.Commit()

	return nil
}

func (db *ConfigDB) GetGlobal(name string) (*GlobalObject, error) {
	t := db.db.Txn(false)
	defer t.Abort()

	e, err := t.First(TableNameGlobal, "id", name)
	if err != nil {
		return nil, err
	}
	if e == nil {
		log.Debugf("Cache MISS on global: %v", name)
		return nil, nil
	}

	log.Debugf("Cache HIT on global: %v", name)

	return e.(*GlobalObject), nil
}

func (db *ConfigDB) DeleteGlobal(name string) error {
	t := db.db.Txn(true)
	defer t.Abort()

	_, err := t.DeleteAll(TableNameGlobal, "id", name)
	if err != nil {
		return err
	}

	t.Commit()

	return nil
}
