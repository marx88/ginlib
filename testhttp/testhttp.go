package testhttp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
)

// Get 根据特定请求uri，发起get请求返回响应
func Get(uri string, router *gin.Engine) []byte {
	// 构造get请求
	req := httptest.NewRequest("GET", uri, nil)
	// 初始化响应
	w := httptest.NewRecorder()

	// 调用相应的handler接口
	router.ServeHTTP(w, req)

	// 提取响应
	result := w.Result()
	defer result.Body.Close()

	// 读取响应body
	body, _ := ioutil.ReadAll(result.Body)
	return body
}

// parseToStr 将map中的键值对输出成querystring形式
func parseToStr(mp map[string]string) string {
	values := []string{}
	for key, val := range mp {
		values = append(values, key+"="+val)
	}
	return strings.Join(values, "&")
}

// PostForm 根据特定请求uri和参数param，以表单形式传递参数，发起post请求返回响应
func PostForm(uri string, param map[string]string, router *gin.Engine) []byte {
	// 构造post请求，表单数据以querystring的形式加在uri之后
	if strings.Index(uri, "?") > -1 {
		uri += parseToStr(param)
	} else {
		uri += "?" + parseToStr(param)
	}
	req := httptest.NewRequest("POST", uri, nil)

	// 初始化响应
	w := httptest.NewRecorder()

	// 调用相应handler接口
	router.ServeHTTP(w, req)

	// 提取响应
	result := w.Result()
	defer result.Body.Close()

	// 读取响应body
	body, _ := ioutil.ReadAll(result.Body)
	return body
}

// PostJSON 根据特定请求uri和参数param，以Json形式传递参数，发起post请求返回响应
func PostJSON(uri string, param map[string]interface{}, router *gin.Engine) []byte {
	// 将参数转化为json比特流
	jsonByte, _ := json.Marshal(param)

	// 构造post请求，json数据以请求body的形式传递
	req := httptest.NewRequest("POST", uri, bytes.NewReader(jsonByte))

	// 初始化响应
	w := httptest.NewRecorder()

	// 调用相应的handler接口
	router.ServeHTTP(w, req)

	// 提取响应
	result := w.Result()
	defer result.Body.Close()

	// 读取响应body
	body, _ := ioutil.ReadAll(result.Body)
	return body
}
