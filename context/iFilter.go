package context

// IFilter 过滤器
type IFilter interface {
	// OnActionExecuting 页面执行前执行
	OnActionExecuting()
	// OnActionExecuted 页面执行后执行
	OnActionExecuted()
}
