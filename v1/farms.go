package v1

import (
	"fmt"

	"github.com/Jeffail/gabs"
	"github.com/gin-gonic/gin"
	"github.com/konsorten/zevenet-rest-api/configdb"
	"github.com/konsorten/zevenet-rest-api/helpers"
	"github.com/konsorten/zevenet-rest-api/models"
)

// GetAllFarms retrieves a list of all farms.
// @Summary Get list of all farms.
// @Description Retrieves a list of all farms. The farms can be up or down.
// @ID get-farms
// @Tags Farms
// @Accept  json
// @Produce  json
// @Success 200 {array} models.FarmInfo
// @Failure 400 {object} models.ApiError
// @Security ApiKeyAuth
// @Router /v1/farms [get]
func (controller *ApiController) GetAllFarms(c *gin.Context) {
	cached, err := configdb.GetInstance().GetGlobal("/farms")
	if err != nil {
		c.AbortWithStatusJSON(400, models.NewApiError(400, "Failed to retrieve from cache: %v", err))
		return
	}

	var content *gabs.Container

	if cached == nil {
		res, err := controller.zapi.Call("GET", "/farms", nil)
		if err != nil {
			c.AbortWithStatusJSON(400, models.NewApiError(400, "Failed to call ZAPI: %v", err))
			return
		}

		content = res.Content

		err = configdb.GetInstance().AddGlobal("/farms", content)
		if err != nil {
			c.AbortWithStatusJSON(400, models.NewApiError(400, "Failed to add to cache: %v", err))
			return
		}
	} else {
		content = cached.Object
	}

	resultBody := content.S("params")

	farms, _ := helpers.GetArray(resultBody)
	for _, farm := range farms {
		farm.Delete("status") // remove uncached value
		farm.Delete("vport")  // remove profile-specific value

		helpers.AddLinks(c, farm, map[string]string{
			"detail": fmt.Sprintf("./%v/%v", farm.S("profile").Data(), farm.S("farmname").Data()),
		})
	}

	c.Data(200, JsonMimeType, resultBody.Bytes())
}
