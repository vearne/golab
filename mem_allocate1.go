package main

import(
	"github.com/imroc/req"
	"fmt"
	"runtime"
)


func init(){
	SetConnPool()
}

func main(){
	size := 1


	var err error
	var resp *req.Resp
	url := "http://up.xiaorui.cc:9000/bigjson.json"
	var mem runtime.MemStats
	for i:=0;i<size + 1;i++{
		resp, err = req.Get(url)
		if err!= nil{
			fmt.Println("error", err)
		}
		fmt.Println(resp.ToString())
		if 1 % 10000 == 0{
			runtime.ReadMemStats(&mem)
			fmt.Println("---------------", i)
			fmt.Println(mem.Alloc)
			fmt.Println(mem.TotalAlloc)
			// 堆对象占用的字节大小
			fmt.Println(mem.HeapAlloc)
			fmt.Println(mem.HeapSys)
			// 分配在堆上的对象数量
			fmt.Println(mem.HeapObjects)
		}
	}

}