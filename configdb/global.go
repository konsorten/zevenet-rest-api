package configdb

import (
	"github.com/Jeffail/gabs"
	"github.com/hashicorp/go-memdb"
)

type GlobalObject struct {
	Name        string
	ConfigFiles []string
	Object      *gabs.Container
}

func getSchemaGlobal() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: "global",
		Indexes: map[string]*memdb.IndexSchema{
			"id": &memdb.IndexSchema{
				Name:    "id",
				Unique:  true,
				Indexer: &memdb.StringFieldIndex{Field: "Name"},
			},
			"cfgFiles": &memdb.IndexSchema{
				Name:    "cfgFiles",
				Indexer: &memdb.StringSliceFieldIndex{Field: "ConfigFiles"},
			},
		},
	}
}

func (db *ConfigDB) AddGlobal(obj *GlobalObject) error {
	t := db.db.Txn(true)
	defer t.Abort()

	err := t.Insert("global", obj)
	if err != nil {
		return err
	}

	t.Commit()

	return nil
}

func (db *ConfigDB) GetGlobal(name string) (*GlobalObject, error) {
	t := db.db.Txn(false)
	defer t.Abort()

	e, err := t.First("global", "id", name)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, nil
	}

	if v, ok := e.(*GlobalObject); ok {
		return v, nil
	}

	return nil, nil
}

func (db *ConfigDB) DeleteGlobal(name string) error {
	t := db.db.Txn(true)
	defer t.Abort()

	_, err := t.DeleteAll("global", "id", name)
	if err != nil {
		return err
	}

	t.Commit()

	return nil
}
