package zapi

import (
	"bytes"
	"fmt"
	"io"
	nlog "log"
	"net/http"
	"net/http/cgi"

	"github.com/Jeffail/gabs"
	"github.com/konsorten/zevenet-rest-api/globalconfig"
	log "github.com/sirupsen/logrus"
)

const (
	ZapiVersion = "3.1"
)

type Result struct {
	HTTPStatus int
	Content    *gabs.Container
	ContentRaw []byte
}

type zapiResponseWriter struct {
	http.ResponseWriter

	data *zapiResponseData
}

type zapiResponseData struct {
	body       bytes.Buffer
	headers    http.Header
	statusCode int
}

func (w zapiResponseWriter) Header() http.Header {
	return w.data.headers
}

func (w zapiResponseWriter) Write(data []byte) (int, error) {
	return w.data.body.Write(data)
}

func (w zapiResponseWriter) WriteHeader(statusCode int) {
	w.data.statusCode = statusCode
}

func Call(method string, path string, input *gabs.Container) (*Result, error) {
	log.Infof("ZAPI call: %v %v", method, path)

	// prepare handler
	handler := &cgi.Handler{
		Path:   fmt.Sprintf("/usr/local/zevenet/www/zapi/v%v/zapi.cgi", ZapiVersion),
		Dir:    "/usr/local/zevenet/www",
		Root:   fmt.Sprintf("/zapi/v%v/zapi.cgi", ZapiVersion),
		Logger: nlog.New(log.New().Writer(), "zapi.cgi", 0),
		Stderr: log.New().Writer(),
	}

	// prepare request
	var reqBody io.Reader
	if input != nil {
		reqBody = bytes.NewReader(input.Bytes())
	} else {
		reqBody = bytes.NewReader(make([]byte, 0))
	}

	req, err := http.NewRequest(method, "https://localhost:444"+handler.Root+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("Failed to create ZAPI request '%v': %v", path, err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("ZAPI_KEY", globalconfig.GetZevenetGlobalConfig().ZAPIKey)

	// perform request
	res := zapiResponseWriter{
		data: &zapiResponseData{
			headers: make(http.Header),
		},
	}

	handler.ServeHTTP(res, req)

	resBodyRaw := res.data.body.Bytes()
	output, _ := gabs.ParseJSON(resBodyRaw)

	// handle result
	return &Result{
		HTTPStatus: res.data.statusCode,
		Content:    output,
		ContentRaw: resBodyRaw,
	}, nil
}
