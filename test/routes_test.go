package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/webapi"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestRoutes(t *testing.T) {
	fs.Initialize[webapi.Module]("demo")
	configure.SetDefault("Log.Component.webapi", true)

	webapi.RegisterRoutes(webapi.Route{Url: "/mini/test1", Method: "POST|GET", Action: func(req pageSizeRequest) string {
		return fmt.Sprintf("hello world pageSize=%d，pageIndex=%d", req.PageSize, req.PageIndex)
	}})
	go webapi.Run(":8084")
	time.Sleep(10 * time.Millisecond)

	t.Run("mini/test1:8084-POST", func(t *testing.T) {
		sizeRequest := pageSizeRequest{PageSize: 10, PageIndex: 2}
		marshal, _ := json.Marshal(sizeRequest)
		rsp, _ := http.Post("http://127.0.0.1:8084/mini/test1", "application/json", bytes.NewReader(marshal))
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "hello world pageSize=10，pageIndex=2", string(body))
		assert.Equal(t, 200, rsp.StatusCode)
	})
	t.Run("mini/test1:8084-GET", func(t *testing.T) {
		rsp, _ := http.Get("http://127.0.0.1:8084/mini/test1?page_size=10&PageIndex=2")
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "hello world pageSize=10，pageIndex=2", string(body))
		assert.Equal(t, 200, rsp.StatusCode)
	})
}
