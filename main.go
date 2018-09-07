package main

import (
	"fmt"
	"net"
	"net/http/fcgi"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/konsorten/zevenet-rest-api/configdb"
	"github.com/konsorten/zevenet-rest-api/globalconfig"
	"github.com/konsorten/zevenet-rest-api/helpers"
	"github.com/konsorten/zevenet-rest-api/v1"
	log "github.com/sirupsen/logrus"
)

const (
	RestApiSocketPath string = "/usr/local/zevenet/app/cherokee/var/run/rest-api.sock"
	LogFilePath       string = "/var/log/rest-api.log"
)

func recoverPanic() {
	if rec := recover(); rec != nil {
		log.Fatalf("FATAL ERROR: %+v", rec)
		os.Exit(1)
	}
}

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
	// setup logger
	logFile, err := os.OpenFile(LogFilePath, os.O_APPEND|os.O_WRONLY|os.O_SYNC, 0644)
	if err != nil {
		return 104, fmt.Errorf("Error creating log file %v: %v", LogFilePath, err)
	}

	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetLevel(log.DebugLevel)

	// dump info
	log.Infof("%v - v%v", mainName, mainVersion)

	defer recoverPanic()

	// start-up components
	err = configdb.CreateConfigDb()
	if err != nil {
		return 121, fmt.Errorf("Error creating Config DB: %v", err)
	}

	err = globalconfig.ReadZevenetGlobalConfig()
	if err != nil {
		return 122, fmt.Errorf("Error reading global config: %v", err)
	}

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
		c.Redirect(301, helpers.ResolveRelativePath(c, "../swagger/"))
	})

	// register v1 root
	_, err = v1.NewApiController(cherokeeRoot)
	if err != nil {
		return 103, fmt.Errorf("Error creating API controller %v: %v", RestApiSocketPath, err)
	}

	// start serving the webserver
	log.Infof("Connecting to FCGI socket: %v", RestApiSocketPath)

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
