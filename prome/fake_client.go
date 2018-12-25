package main

import (
	"flag"
	"fmt"
	"github.com/imroc/req"
	"math/rand"
	"net/http"
	"time"
)

func SetConnPool() {
	client := &http.Client{}
	client.Transport = &http.Transport{
		MaxIdleConnsPerHost: 500,
		// 无需设置MaxIdleConns
		// MaxIdleConns controls the maximum number of idle (keep-alive)
		// connections across all hosts. Zero means no limit.
		// MaxIdleConns 默认是0，0表示不限制
	}

	req.SetClient(client)
	req.SetTimeout(5 * time.Second)
}

func main() {
	SetConnPool()
	host := "127.0.0.1"
	flag.StringVar(&host, "host", "127.0.0.1", "127.0.0.1")

	flag.Parse()
	url1 := fmt.Sprintf("http://%s:28181/api/api1", host)
	url2 := fmt.Sprintf("http://%s:28181/api/api2", host)

	fmt.Println("url1", url1, "url2", url2)

	// use Req object to initiate requests.
	var count = 0
	r := req.New()

	c := time.Tick(time.Millisecond * 100)
	for {
		select {
		case <-c:
			count = rand.Intn(20)
			for i := 0; i < count; i++ {
				resp, _ := r.Get(url1)
				resp.Bytes()
			}

			for i := 0; i < count*2; i++ {
				resp, _ := r.Get(url2)
				resp.Bytes()
			}
		}
	}
}
