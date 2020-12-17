package router

import "github.com/gin-gonic/gin"

// RegisterFunc 路由注册函数类型
type RegisterFunc func(*gin.RouterGroup)

// Group 路由组
type Group struct {
	path        string            // 路由组路径
	middlewares gin.HandlersChain // 中间件列表
	registers   []RegisterFunc    // 路由注册函数列表
	subGroups   []*Group          // 子路由组列表
}

// NewGroup 实例化路由组
func NewGroup(path string) *Group {
	return &Group{
		path,
		make(gin.HandlersChain, 0),
		[]RegisterFunc{},
		[]*Group{},
	}
}

// AddMiddleware 添加中间件
func (g *Group) AddMiddleware(middleware ...gin.HandlerFunc) {
	g.middlewares = append(g.middlewares, middleware...)
}

// AddRegister 添加路由注册函数
func (g *Group) AddRegister(register ...RegisterFunc) {
	g.registers = append(g.registers, register...)
}

// AddSubGroup 添加子路由组
func (g *Group) AddSubGroup(subGroup ...*Group) {
	g.subGroups = append(g.subGroups, subGroup...)
}

// Mount 挂载路由组
func (g *Group) Mount(parentRouter gin.IRouter) {
	// 挂载路由组到父路由上
	routerGroup := parentRouter.Group(g.path)

	// 挂载中间件
	routerGroup.Use(g.middlewares...)

	// 注册子路由
	g.execRegister(routerGroup)

	// 挂载子路由组
	g.mountSubGroup(routerGroup)
}

// execRegister 注册子路由
func (g *Group) execRegister(routerGroup *gin.RouterGroup) {
	for _, register := range g.registers {
		register(routerGroup)
	}
}

// mountSubGroup 挂载子路由组
func (g *Group) mountSubGroup(parentRouterGroup *gin.RouterGroup) {
	for _, subGroup := range g.subGroups {
		subGroup.Mount(parentRouterGroup)
	}
}
