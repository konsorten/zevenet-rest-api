package main

import (
	"fmt"
	"net"
	"net/http/fcgi"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/konsorten/zevenet-rest-api/v1"
	log "github.com/sirupsen/logrus"
)

const (
	RestApiSocketPath string = "/usr/local/zevenet/app/cherokee/var/run/rest-api.sock"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /rest-api
func main() {
	ret, err := mainInternal()
	if err != nil {
		log.Errorf("ERROR: %v", err)
		os.Exit(ret)
		return
	}

	if ret != 0 {
		os.Exit(ret)
	}
}

func mainInternal() (int, error) {
	// register the handlers
	handler := gin.New()

	handler.Use(func(c *gin.Context) {
		// inject missing values
		c.Request.RequestURI = c.Request.URL.String()
		c.Next()
	})

	handler.Use(gin.Recovery())

	// register /rest-api/
	cherokeeRoot := handler.Group("/rest-api")

	cherokeeRoot.GET("/", func(c *gin.Context) {
		c.Redirect(301, "../swagger/")
	})

	// register v1 root
	v1Controller, err := v1.NewApiController(cherokeeRoot)
	if err != nil {
		return 103, fmt.Errorf("Error creating API controller %v: %v", RestApiSocketPath, err)
	}

	v1Controller.Register()

	// start serving the webserver
	os.Remove(RestApiSocketPath) // try to remove unclean socket

	listener, err := net.Listen("unix", RestApiSocketPath)
	if err != nil {
		return 101, fmt.Errorf("Error listening to %v: %v", RestApiSocketPath, err)
	}

	defer listener.Close()

	// start serving requests
	err = fcgi.Serve(listener, handler)
	if err != nil {
		return 102, fmt.Errorf("Error serving: %v", err)
	}

	return 0, nil
}
