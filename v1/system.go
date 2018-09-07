package v1

import (
	"github.com/Jeffail/gabs"
	"github.com/gin-gonic/gin"
	"github.com/konsorten/zevenet-rest-api/configdb"
	"github.com/konsorten/zevenet-rest-api/models"
	"github.com/konsorten/zevenet-rest-api/zapi"
)

// GetSystemVersion retrieves basic system information like various version numbers.
// @Summary Get basic system information.
// @Description Retrieves basic system information like various version numbers.
// @ID get-system-version
// @Tags System
// @Accept  json
// @Produce  json
// @Success 200 {object} models.SystemVersion
// @Failure 400 {object} models.ApiError
// @Security ApiKeyAuth
// @Router /v1/system/version [get]
func (controller *ApiController) GetSystemVersion(c *gin.Context) {
	cached, err := configdb.GetInstance().GetGlobal("/system/version")
	if err != nil {
		c.AbortWithStatusJSON(400, models.NewApiError(400, "Failed to retrieve from cache: %v", err))
		return
	}

	var content *gabs.Container

	if cached == nil {
		res, err := zapi.Call("GET", "/system/version", nil)
		if err != nil {
			c.AbortWithStatusJSON(400, models.NewApiError(400, "Failed to call ZAPI: %v", err))
			return
		}

		content = res.Content

		err = configdb.GetInstance().AddGlobal("/system/version", content)
		if err != nil {
			c.AbortWithStatusJSON(400, models.NewApiError(400, "Failed to add to cache: %v", err))
			return
		}
	} else {
		content = cached.Object
	}

	resultBody := content.S("params")

	resultBody.Set(mainVersionSimple, "rest_api_version")

	c.Data(200, JsonMimeType, resultBody.Bytes())
}
