package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	RestApiSocketPath string = "/usr/local/zevenet/app/cherokee/var/run/rest-api.sock"
)

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
	handler := http.NewServeMux()

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})

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
