package middleware

import (
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/webapi/context"
	"net/http"
	"reflect"
)

type ApiResponse struct {
	context.IMiddleware
}

func (receiver *ApiResponse) Invoke(httpContext *context.HttpContext) {
	// ActionResult类型，不做ApiResponse解析
	if httpContext.IsActionResult() {
		receiver.IMiddleware.Invoke(httpContext)
		return
	}

	var apiResponse core.ApiResponse[any]
	exception.Try(func() {
		receiver.IMiddleware.Invoke(httpContext)

		var returnVal any
		// 只有一个返回值
		if len(httpContext.Response.Body) == 1 {
			returnVal = httpContext.Response.Body[0].Interface()
		} else {
			// 多个返回值，则转成数组Json
			lst := collections.NewListAny()
			for i := 0; i < len(httpContext.Response.Body); i++ {
				lst.Add(httpContext.Response.Body[i].Interface())
			}
			returnVal = lst
		}
		apiResponse = core.Success[any](httpContext.Response.GetStatusMessage(), returnVal)
		apiResponse.StatusCode = httpContext.Response.GetStatusCode()
	}).CatchWebException(func(exp exception.WebException) {
		// 响应码
		httpContext.Exception = exp.Message
		apiResponse = core.Error[any](exp.Message, exp.StatusCode)
	}).CatchException(func(exp any) {
		// 响应码
		httpContext.Exception = exp
		apiResponse = core.Error[any](fmt.Sprint(exp), http.StatusInternalServerError)
	})

	httpContext.Route.IsGoBasicType = false
	httpContext.Response.Body = []reflect.Value{reflect.ValueOf(apiResponse)}
}
