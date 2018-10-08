package main

import(
	"github.com/imroc/req"
	"fmt"
	"runtime"
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

func init(){
	SetConnPool()
}

func main(){
	size := 0

	var err error
	var resp *req.Resp
	url := "http://up.xiaorui.cc:9000/bigjson.json"
	var mem runtime.MemStats
	for i:=0;i<size + 1;i++{
		resp, err = req.Get(url)
		if err!= nil{
			fmt.Println("error", err)
			continue
		}
		data,  _ := resp.ToBytes()
		fmt.Println("resp byte length", len(data))
		if i % 10000 == 0{
			runtime.ReadMemStats(&mem)
			fmt.Println("---------------", i)
			fmt.Println("mem.Alloc", mem.Alloc)
			fmt.Println("mem.TotalAlloc", mem.TotalAlloc)
			// 堆对象占用的字节大小
			fmt.Println("mem.HeapAlloc", mem.HeapAlloc)
			fmt.Println("mem.HeapSys", mem.HeapSys)
			// 分配在堆上的对象数量
			fmt.Println("mem.HeapObjects", mem.HeapObjects)
		}
	}

}