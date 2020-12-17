package router

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/marx88/ginlib/testhttp"
)

// 简单的测试

type testUserController struct{}

var (
	testEngine     *gin.Engine
	testAdminGroup *Group
)

func (uc *testUserController) index(c *gin.Context) {
	c.String(http.StatusOK, "user index")
}

func (uc *testUserController) use(c *gin.Context) {
	before, after := c.GetString("before"), c.GetString("after")
	c.String(http.StatusOK, before+after)
}

// TestSuite 测试套件
func TestSuite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testEngine = gin.Default()
	t.Run("测试 新建路由组", testNewGroup)
	t.Run("测试 添加中间件", testAddMiddleware)
	t.Run("测试 添加路由注册函数", testAddRegister)
	t.Run("测试 添加子路由组", testAddSubGroup)
	t.Run("测试 挂载路由组", testMount)
}

func testNewGroup(t *testing.T) {
	if testAdminGroup = NewGroup("admin"); testAdminGroup == nil {
		t.Error("func NewGroup failed.")
	}
}

func testAddMiddleware(t *testing.T) {
	expected := len(testAdminGroup.middlewares) + 1
	testAdminGroup.AddMiddleware(func(c *gin.Context) {})
	if result := len(testAdminGroup.middlewares); result != expected {
		t.Errorf(
			"func AddMiddleware failed, expected: %d, result: %d.",
			expected,
			result,
		)
	}
}

func testAddRegister(t *testing.T) {
	expected := len(testAdminGroup.registers) + 1
	testAdminGroup.AddRegister(func(rg *gin.RouterGroup) {})
	if result := len(testAdminGroup.registers); result != expected {
		t.Errorf(
			"func AddRegister failed, expected: %d, result: %d.",
			expected,
			result,
		)
	}
}

func testAddSubGroup(t *testing.T) {
	expected := len(testAdminGroup.subGroups) + 1
	testAdminGroup.AddSubGroup(NewGroup("user"))
	if result := len(testAdminGroup.subGroups); result != expected {
		t.Errorf(
			"func AddSubGroup failed, expected: %d, result: %d.",
			expected,
			result,
		)
	}
}

func testMount(t *testing.T) {
	// admin路由组
	testAdminGroup = NewGroup("admin")
	testAdminGroup.AddMiddleware(func(c *gin.Context) {
		c.Set("before", "before")
		c.Next()
		c.Set("after", "after")
	})

	// user路由组
	userGroup := NewGroup("user")
	userGroup.AddRegister(func(rg *gin.RouterGroup) {
		uc := new(testUserController)
		rg.GET("index", uc.index)
		rg.POST("use", uc.use)
	})
	testAdminGroup.AddSubGroup(userGroup)

	testAdminGroup.Mount(testEngine)

	// 测试 /admin/user/index 路由
	if resp := string(testhttp.Get("/admin/user/index", testEngine)); resp != "user index" {
		t.Errorf("func Mount failed, router {GET /admin/user/index} failed, expected: user index, resp: %s.", resp)
	}

	// 测试中间件
	params := make(map[string]string, 0)
	if resp := string(testhttp.PostForm("/admin/user/use", params, testEngine)); resp != "before" {
		t.Errorf("func Mount failed, router {POST /admin/user/use} failed, expected: before, resp: %s,", resp)
	}
}
