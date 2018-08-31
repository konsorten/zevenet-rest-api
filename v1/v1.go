package v1

import (
	"github.com/gin-gonic/gin"
)

type ApiController struct {
	handler *gin.RouterGroup
}

func NewApiController(handler *gin.RouterGroup) (*ApiController, error) {
	return &ApiController{
		handler: handler,
	}, nil
}
