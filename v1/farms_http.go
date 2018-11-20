package v1

import (
	"fmt"
	"strconv"

	"github.com/Jeffail/gabs"
	"github.com/gin-gonic/gin"
	"github.com/konsorten/zevenet-rest-api/models"
	"github.com/konsorten/zevenet-rest-api/zapi"
)

type HTTPFarmApiController struct {
	farmSpecificController

	farmType models.FarmProfile
	generic  *farmApiController
	handler  *gin.RouterGroup
	zapi     *zapi.Host
}

func NewHTTPFarmApiController(handler *gin.RouterGroup) (*HTTPFarmApiController, error) {
	controller := &HTTPFarmApiController{
		handler:  handler,
		zapi:     zapi.NewHost("3.1"),
		farmType: models.FarmProfile_HTTP,
	}

	genericController, err := newFarmApiController(handler, controller)
	if err != nil {
		return nil, err
	}
	controller.generic = genericController

	handler.GET("", controller.GetAllFarms)

	return controller, nil
}

func (controller *HTTPFarmApiController) GetFarmProfile() models.FarmProfile {
	return controller.farmType
}

func (controller *HTTPFarmApiController) GetConfigFileList(farmName string) []string {
	return []string{
		fmt.Sprintf("%v_pound.cfg", farmName),
	}
}

func (controller *HTTPFarmApiController) parseVirtualPort(parent *gabs.Container, path string) error {
	vport := parent.S(path).Data()
	if s, ok := vport.(string); ok {
		i, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		parent.Set(i, path)
		return nil
	}

	return fmt.Errorf("Failed to convert virtual port: %+v", vport)
}

// GetAllFarms retrieves a list of all HTTP farms.
// @Summary Get list of all HTTP farms.
// @Description Retrieves a list of all HTTP farms. The farms can be up or down. Does not include HTTPS farms.
// @ID get-farms-http
// @Tags Farms
// @Accept  json
// @Produce  json
// @Success 200 {array} models.FarmInfo
// @Failure 400 {object} models.ApiError
// @Security ApiKeyAuth
// @Router /v1/farms/http [get]
func (controller *HTTPFarmApiController) GetAllFarms(c *gin.Context) {
	c.Header(SwaggerIDHeader, "get-farms-http")

	controller.generic.GetAllFarms(c)
}
