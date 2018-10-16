package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"net/http"
	"bytes"
	"encoding/json"
)


const SUCCESS_CODE = "E000"
const INTERNAL_CODE = "E001"

type ErrorResp struct {
	Code string `json:"code"`
	Msg string `json:"msg"`
}

type ErrorCode struct{
	Code string `json:"code"`
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	if n, err := w.body.Write(b); err != nil {
		log.Println("err", err)
		return n, err
	}
	return w.ResponseWriter.Write(b)
}

func Logger2() gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request
		// 注意这里Buff的创建可能带来性能问题
		// 可参考 [GOLANG 内存分配优化](http://vearne.cc/archives/671)
		// 进行优化
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		var code string
		var res ErrorCode

		log.Println(string(blw.body.Bytes()))
		if err := json.Unmarshal(blw.body.Bytes(), &res); err != nil {
			// Body体不是合法的JSON或"code"字段类型不一致
			code = INTERNAL_CODE
			log.Println("error", err)
		} else if len(res.Code) <= 0 {
			// Body体是合法JSON
			// JSON不含有"code"字段
			code = SUCCESS_CODE
		}else{
			// Body体是合法JSON
			// JSON含有"code"字段
			code = res.Code
		}

		log.Println("code", code)
	}
}


func main() {
	r := gin.New()
	r.Use(Logger2())

	r.GET("/test", func(c *gin.Context) {
		log.Println("-----gin_middleware2------")
		//data := `{"msg":"test"}`
		//data := `xxx`
		//data := `{"code":"E002", "msg":"test"}`
		data := `{"code": 10, "msg":"test"}`
		c.Data(http.StatusOK,
			"application/json", []byte(data))
	})

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}