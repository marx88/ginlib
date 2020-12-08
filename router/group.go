package router

import "github.com/gin-gonic/gin"

// Group 路由组
type Group struct {
	path              string             // 路由组路径
	middlewares       []gin.HandlerFunc  // 中间件列表
	subRouterHandlers []GroupHandlerFunc // 子路由注册函数列表
	subRouterGroups   []*Group           // 子路由组列表
}

// GroupHandlerFunc 子路由注册函数类型
type GroupHandlerFunc func(*gin.RouterGroup)

// NewGroup 实例化一个RouterGroup指针
func NewGroup(path string) *Group {
	return &Group{
		path,
		[]gin.HandlerFunc{},
		[]GroupHandlerFunc{},
		[]*Group{},
	}
}

// AddMiddleware 添加中间件
func (r *Group) AddMiddleware(middleware ...gin.HandlerFunc) {
	r.middlewares = append(r.middlewares, middleware...)
}

// AddSubRouter 添加子路由注册函数
func (r *Group) AddSubRouter(subRouterHandler ...GroupHandlerFunc) {
	r.subRouterHandlers = append(r.subRouterHandlers, subRouterHandler...)
}

// AddSubRouterGroup 添加子路由组
func (r *Group) AddSubRouterGroup(subRouterGroup ...*Group) {
	r.subRouterGroups = append(r.subRouterGroups, subRouterGroup...)
}

// Mount 挂载路由组
func (r *Group) Mount(parentRouter gin.IRouter) {
	// 挂载路由组到父路由上
	routerGroup := parentRouter.Group(r.path)

	// 挂载中间件
	routerGroup.Use(r.middlewares...)

	// 挂载子路由
	r.mountSubRouter(routerGroup)

	// 挂载子路由组
	r.mountSubRouterGroup(routerGroup)
}

// mountSubRouter 挂载子路由
func (r *Group) mountSubRouter(routerGroup *gin.RouterGroup) {
	for _, handler := range r.subRouterHandlers {
		handler(routerGroup)
	}
}

// mountSubRouterGroup 挂载子路由组
func (r *Group) mountSubRouterGroup(parentRouterGroup *gin.RouterGroup) {
	for _, subRouterGroup := range r.subRouterGroups {
		subRouterGroup.Mount(parentRouterGroup)
	}
}
