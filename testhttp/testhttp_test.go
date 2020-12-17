package testhttp

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

type testPostFormType struct {
	Name string `form:"name" json:"name" binding:"required"`
	Age  int32  `form:"age" json:"age" binding:"required"`
}

type testRespType struct {
	Code int32            `json:"code"`
	Msg  string           `json:"msg"`
	Data testPostFormType `json:"data"`
}

var testRouter *gin.Engine

func init() {
	gin.SetMode(gin.TestMode)
	testRouter = gin.New()
	testRouter.Use(gin.Recovery())
	testRouter.GET("/get/:name", func(c *gin.Context) {
		c.String(http.StatusOK, c.Param("name"))
	})
	testRouter.POST("/post", func(c *gin.Context) {
		form := new(testPostFormType)
		if err := c.ShouldBind(form); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "error",
				"data": nil,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "success",
			"data": form,
		})
	})
	testRouter.POST("/json", func(c *gin.Context) {
		form := new(testPostFormType)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "error",
				"data": nil,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "success",
			"data": form,
		})
	})
}

func TestGet(t *testing.T) {
	if resp := string(Get("/get/test", testRouter)); resp != "test" {
		t.Error("func Get failed，expected：test，resp：" + resp)
	}
}

func TestPostForm(t *testing.T) {
	params := make(map[string]string)
	params["name"] = "test"
	params["age"] = "18"

	resp := new(testRespType)
	if err := json.Unmarshal(PostForm("/post", params, testRouter), resp); err != nil {
		t.Errorf("func PostForm error：%v", err)
	}
	if resp.Data.Name != "test" {
		t.Errorf("func PostForm failed，expected：test，result：%s", resp.Data.Name)
	}
}

func TestPostJSON(t *testing.T) {
	params := make(map[string]interface{})
	params["name"] = "test"
	params["age"] = 18

	resp := new(testRespType)
	if err := json.Unmarshal(PostJSON("/json", params, testRouter), resp); err != nil {
		t.Errorf("func PostJSON error：%v", err)
	}
	if resp.Data.Name != "test" {
		t.Errorf("func PostJSON failed，expected：test，result：%s", resp.Data.Name)
	}
}
