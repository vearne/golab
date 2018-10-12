package main

import (
	"github.com/imroc/req"
	"fmt"
	"runtime"
	"sync"
	"bytes"
	"net/http"
	"time"
)

func SetConnPool() {
	client := &http.Client{}
	client.Transport = &http.Transport{
		MaxIdleConnsPerHost: 1000,
		// 无需设置MaxIdleConns
		// MaxIdleConns controls the maximum number of idle (keep-alive)
		// connections across all hosts. Zero means no limit.
		// MaxIdleConns 默认是0，0表示不限制
	}

	req.SetClient(client)
	req.SetTimeout(5 * time.Second)
}

var buffPool sync.Pool

func init() {
	SetConnPool()
}

func main() {
	size := 10000

	var err error
	var resp *req.Resp
	var buffer *bytes.Buffer
	var byteSlice []byte
	url := "http://localhost:9000/bigjson.json"
	var mem runtime.MemStats
	for i := 0; i < size+1; i++ {
		resp, err = req.Get(url)
		if err != nil {
			fmt.Println("error", err)
			continue
		}

		item := buffPool.Get()
		if item == nil {
			byteSlice = make([]byte, 0, 10*1024)
		} else {
			byteSlice = item.([]byte)
		}
		buffer = bytes.NewBuffer(byteSlice)
		buffer.ReadFrom(resp.Response().Body)


		res  := buffer.Bytes()
		fmt.Println("resp byte length", len(res))

		//fmt.Println("len", buffer.Len())
		//fmt.Println("Cap", buffer.Cap())
		buffPool.Put(byteSlice)

		if i%1000 == 0 {
			runtime.ReadMemStats(&mem)
			fmt.Println("---------------", i)
			fmt.Println("mem.Alloc", mem.Alloc)
			// 为堆对象总计分配的字节数
			fmt.Println("mem.TotalAlloc", mem.TotalAlloc)
			// 为创建堆对象总计的内存申请次数
			fmt.Println("mem.Mallocs", mem.Mallocs)
			// 为销毁堆对象总计的内存释放次数
			fmt.Println("mem.Frees", mem.Frees)
			// 堆对象占用的字节大小
			fmt.Println("mem.HeapAlloc", mem.HeapAlloc)
			fmt.Println("mem.HeapSys", mem.HeapSys)
			// 分配在堆上的对象数量
			fmt.Println("mem.HeapObjects", mem.HeapObjects)
		}
	}

}
