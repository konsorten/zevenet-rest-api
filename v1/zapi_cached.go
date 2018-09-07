package v1

import (
	"github.com/Jeffail/gabs"
	"github.com/konsorten/zevenet-rest-api/configdb"
	"github.com/konsorten/zevenet-rest-api/models"
	"github.com/konsorten/zevenet-rest-api/zapi"
)

func callZAPICached(zapiHost *zapi.Host, getPath string) (*gabs.Container, error) {
	cached, err := configdb.GetInstance().GetGlobal(getPath)
	if err != nil {
		return nil, models.NewApiError(400, "Failed to retrieve from cache: %v", err)
	}

	var content *gabs.Container

	if cached == nil {
		res, err := zapiHost.Call("GET", getPath, nil)
		if err != nil {
			return nil, models.NewApiError(400, "Failed to call ZAPI: %v", err)
		}

		content = res.Content

		err = configdb.GetInstance().AddGlobal(getPath, content)
		if err != nil {
			return nil, models.NewApiError(400, "Failed to add to cache: %v", err)
		}
	} else {
		content = cached.Object()
	}

	return content.S("params"), nil
}
