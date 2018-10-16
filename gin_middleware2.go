package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"net/http"
	"bytes"
	"encoding/json"
)


const SUCCESS_CODE = "E000"

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
		if err := json.Unmarshal(blw.body.Bytes(), &res); err != nil {
			code = res.Code
		} else {
			code = SUCCESS_CODE
		}
		log.Println("code", code)
	}
}


func main() {
	r := gin.New()
	r.Use(Logger2())

	r.GET("/test", func(c *gin.Context) {
		log.Println("-----gin_middleware2------")
		c.JSON(http.StatusOK, ErrorResp{Code:"E001", Msg:"params error"})
	})

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}