package api

import (
	service2 "SService/service"
)

// HandlerGroup 包含所有处理器的结构
type ApiGroup struct {
	ExampleApi
	UserApi
	HTMLRecordApi
	HTMLRecordDataApi
}

var (
	exampleService        = service2.ExampleService{}
	userService           = service2.UserService{}
	htmlRecordService     = service2.HTMLRecordService{}
	htmlRecordDataService = service2.HTMLRecordDataService{}
)
