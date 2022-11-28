package minimal

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/types"
	"os"
	"reflect"
)

// Register 注册单个Api
func Register(area string, method string, route string, actionFunc any, paramNames ...string) {
	actionType := reflect.TypeOf(actionFunc)
	param := types.GetInParam(actionType)

	// 如果设置了方法的入参（多参数），则需要全部设置
	if len(paramNames) > 0 && len(paramNames) != len(param) {
		flog.Errorf("注册minimalApi失败：%s函数入参个数设置与%s不匹配", flog.Colors[eumLogLevel.Error](actionType.String()), flog.Colors[eumLogLevel.Error](paramNames))
		os.Exit(1)
	}

	// 添加到路由表
	lstRouteTable.Add(routeTable{
		routeUrl:         area + route,
		action:           actionFunc,
		method:           method,
		requestParamType: collections.NewList(param...),
		responseBodyType: collections.NewList(types.GetOutParam(actionType)...),
		paramNames:       collections.NewList(paramNames...),
	})
}
