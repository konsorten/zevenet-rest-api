package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/konsorten/zevenet-rest-api/globalconfig"
	"github.com/konsorten/zevenet-rest-api/models"
)

const (
	JsonMimeType = "application/json"
)

type ApiController struct {
	handler *gin.RouterGroup
}

func NewApiController(handler *gin.RouterGroup) (*ApiController, error) {
	// read the global configuration
	controller := &ApiController{
		handler: handler,
	}

	// setup authentication
	handler.Use(func(c *gin.Context) {
		globalConfig := globalconfig.GetZevenetGlobalConfig()

		if globalConfig.ZAPIKey == "" {
			c.AbortWithStatusJSON(403, models.NewApiError(403, "ZAPI user is disabled or no ZAPI key set"))
			return
		}
		if c.GetHeader("Zapi-Key") != globalConfig.ZAPIKey {
			c.AbortWithStatusJSON(401, models.NewApiError(401, "ZAPI key does not match (see ZAPI_KEY header)"))
			return
		}
		c.Next()
	})

	// register endpoints
	v1 := handler.Group("/v1")

	v1.GET("/", controller.GetRoot)
	v1.GET("/system/version", controller.GetSystemVersion)

	return controller, nil
}
