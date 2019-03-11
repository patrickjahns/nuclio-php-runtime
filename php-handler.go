/**
 * @author Patrick Jahns <github@patrickjahns.de>
 *
 * @copyright Copyright (c) 2019, Patrick Jahns.
 * @license Apache
 *
 */

package main

import (
	"bytes"
	"github.com/nuclio/nuclio-sdk-go"
	"github.com/tomasen/fcgi_client"
	"github.com/yookoala/gofast/tools/phpfpm"
	"io/ioutil"
	"os"
	"strings"
)

func Handler(context *nuclio.Context, event nuclio.Event) (interface{}, error) {
	context.Logger.Info("invocation started")

	fpmSocket := "/var/task/fpm.sock"
	fpmConfFile := "/var/task/php-fpm.conf"
	fpm := phpfpm.NewProcess(os.Getenv("PHP_FPM_BIN"))
	fpm.SetDatadir("/var/task")
	fpm.ConfigFile = fpmConfFile
	fpm.Listen = fpmSocket

	if _, err := os.Stat(fpmSocket); os.IsNotExist(err) {
		context.Logger.Info("starting fpm")
		err := fpm.Start()
		if err != nil {
			context.Logger.Error("err:", err)
			return nuclio.Response{StatusCode: 500}, nil
		}
	}

	context.Logger.Info("connecting to fcgi")
	fcgi, err := fcgiclient.Dial("unix", fpmSocket)
	if err != nil {
		context.Logger.Error("err:", err)
		return nuclio.Response{StatusCode: 500}, nil

	}
	defer fcgi.Close()

	scriptFilename := os.Getenv("PHP_SCRIPT")
	env := map[string]string{
		"CONTENT_LENGTH":  event.GetHeaderString(string("Content-Length")),
		"CONTENT_TYPE":    event.GetContentType(),
		"REQUEST_METHOD":  event.GetMethod(),
		"PATH_INFO":       event.GetPath(),
		"REQUEST_URI":     "/" + event.GetHeaderString("X-Nuclio-Path"),
		"SCRIPT_FILENAME": scriptFilename,
		"SERVER_SOFTWARE": "nuclio-go / fcgiclient ",
	}

	for header, _ := range event.GetHeaders() {
		env["HTTP_"+strings.Replace(strings.ToUpper(header), "-", "_", -1)] = event.GetHeaderString(header)
	}

	reqBody := bytes.NewReader(event.GetBody())

	context.Logger.Info("sending request to fpm", env)
	resp, err := fcgi.Request(env, reqBody)
	if err != nil {
		context.Logger.Error("err:", err)
		return nuclio.Response{StatusCode: 500}, nil
	}

	respHeaders := make(map[string]interface{})
	for k, _ := range resp.Header {
		respHeaders[k] = resp.Header.Get(k)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		context.Logger.Error("err:", err)
		return nuclio.Response{StatusCode: 500}, nil
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if len(contentType) == 0 {
	    contentType = "text/plain"
	}

	return nuclio.Response{
		StatusCode:  resp.StatusCode,
		ContentType: contentType,
		Body:        []byte(content),
		Headers:     respHeaders,
	}, nil

}
