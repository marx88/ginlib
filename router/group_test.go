package router

import (
	"testing"

	"github.com/gin-gonic/gin"
)

// 简单的测试

// TestNewGroup 测试 新建路由组
func TestNewGroup(t *testing.T) {
	NewGroup("/")
}

// TestAddMiddleware 测试 给路由组添加中间件
func TestAddMiddleware(t *testing.T) {
	group := NewGroup("/")
	group.AddMiddleware(func(c *gin.Context) {})
}

// TestAddSubRouter 测试 给路由组添加子路由
func TestAddSubRouter(t *testing.T) {
	group := NewGroup("/")
	group.AddSubRouter(func(r *gin.RouterGroup) {})
}

// TestAddSubRouterGroup 测试 给路由组添加子路由组
func TestAddSubRouterGroup(t *testing.T) {
	group := NewGroup("/")
	group.AddSubRouterGroup(NewGroup("/admin"))
}

// TestMount 测试 挂载路由组到gin.IRouter
func TestMount(t *testing.T) {
	// 创建根路由组及其中间件、子路由
	root := NewGroup("/")
	root.AddMiddleware(func(c *gin.Context) {})
	root.AddSubRouter(func(r *gin.RouterGroup) {})

	// 创建admin路由组及其中间件、子路由
	admin := NewGroup("/admin")
	admin.AddMiddleware(func(c *gin.Context) {})
	admin.AddSubRouter(func(r *gin.RouterGroup) {})

	// admin路由组挂载到root路由组下
	root.AddSubRouterGroup(admin)

	// root路由组挂载到gin.engine下
	root.Mount(gin.New())
}
