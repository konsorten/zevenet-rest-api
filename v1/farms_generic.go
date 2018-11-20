package v1

import (
	"fmt"

	"github.com/Jeffail/gabs"
	"github.com/gin-gonic/gin"
	"github.com/konsorten/zevenet-rest-api/helpers"
	"github.com/konsorten/zevenet-rest-api/models"
	"github.com/konsorten/zevenet-rest-api/zapi"
)

type farmSpecificController interface {
	GetFarmProfile() models.FarmProfile
	GetConfigFileList(farmName string) []string

	parseVirtualPort(parent *gabs.Container, path string) error
}

type farmApiController struct {
	handler  *gin.RouterGroup
	specific farmSpecificController
	zapi     *zapi.Host
}

func newFarmApiController(handler *gin.RouterGroup, specific farmSpecificController) (*farmApiController, error) {
	controller := &farmApiController{
		handler:  handler,
		specific: specific,
		zapi:     zapi.NewHost("3.1"),
	}

	// do net register any handlers
	// they are registered by the specific controllers
	// like HTTPFarmApiController

	return controller, nil
}

func (controller *farmApiController) GetAllFarms(c *gin.Context) {
	resultBody, err := callZAPICached(controller.zapi, "/farms")
	if err != nil {
		if apiErr, ok := err.(models.ApiError); ok {
			c.AbortWithStatusJSON(apiErr.StatusCode, apiErr)
			return
		} else {
			c.AbortWithStatusJSON(400, models.NewApiError(400, "Unexpected error: %v", err))
			return
		}
	}

	farms, _ := resultBody.Children()
	farmProfile := controller.specific.GetFarmProfile()
	for i, farm := range farms {
		profile := farm.S("profile").Data()
		if string(farmProfile) != profile {
			resultBody.ArrayRemove(i)
			continue
		}

		farm.Delete("status") // remove uncached value

		err := controller.specific.parseVirtualPort(farm, "vport")
		if err != nil {
			c.AbortWithStatusJSON(400, models.NewApiError(400, "Failed to parse virtual port: %v", err))
			return
		}

		helpers.AddLinks(c, farm, map[string]string{
			"detail": fmt.Sprintf("./%v/%v", profile, farm.S("farmname").Data()),
		})
	}

	c.Data(200, JsonMimeType, resultBody.Bytes())
}
