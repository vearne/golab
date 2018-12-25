package main

import (
	"flag"
	"fmt"
	"github.com/imroc/req"
	"math/rand"
	"time"
)

func main() {
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
				r.Get(url1)
			}

			for i := 0; i < count*2; i++ {
				r.Get(url2)
			}
		}
	}
}
