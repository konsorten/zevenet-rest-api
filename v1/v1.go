package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	JsonMimeType = "application/json"
)

type ApiController struct {
	handler      *gin.RouterGroup
	globalConfig *ZevenetGlobalConfig
	lastError    error
}

func NewApiController(handler *gin.RouterGroup) (*ApiController, error) {
	// read the global configuration
	globalConfig, err := ReadZevenetGlobalConfig()
	if err != nil {
		return nil, fmt.Errorf("Error reading Zevenet global config: %v", err)
	}

	controller := &ApiController{
		handler:      handler,
		globalConfig: globalConfig,
	}

	// setup authentication
	handler.Use(func(c *gin.Context) {
		if globalConfig.ZAPIKey == "" {
			c.AbortWithStatusJSON(403, NewApiError(403, "ZAPI user is disabled or no ZAPI key set"))
			return
		}
		if c.GetHeader("Zapi-Key") != globalConfig.ZAPIKey {
			c.AbortWithStatusJSON(401, NewApiError(401, "ZAPI key does not match (see ZAPI_KEY header)"))
			return
		}
		c.Next()
	})

	// register endpoints
	v1 := handler.Group("/v1")

	v1.GET("/", controller.GetRoot)

	// handle errors
	handler.Use(func(c *gin.Context) {
		if controller.lastError == nil {
			c.Next()
			return
		}

		// handle error type
		switch e := controller.lastError.(type) {
		case *ApiError:
			c.AbortWithStatusJSON(e.StatusCode, e)
		default:
			c.AbortWithStatusJSON(500, NewApiError(500, e.Error()))
		}
	})

	return controller, nil
}

func (controller *ApiController) Fail(err error) {
	controller.lastError = err
}
