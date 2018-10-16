package main

import (
	"fmt"
	"bytes"
	"github.com/xtaci/kcp-go"

)

const (
	mutexLocked = 1 << iota // mutex is locked
	mutexWoken
	mutexStarving
	mutexWaiterShift      = iota
	starvationThresholdNs = 1e6
)

func main() {
	fmt.Println("mutexLocked", mutexLocked)
	fmt.Println("mutexWoken", mutexWoken)
	fmt.Println("mutexStarving", mutexStarving)
	fmt.Println("mutexWaiterShift", mutexWaiterShift)
	fmt.Println("starvationThresholdNs", starvationThresholdNs)

	var buff bytes.Buffer
	buff.Grow(1000)
	fmt.Println("cap---1", buff.Cap())
	buff.Reset()
	fmt.Println("cap---2", buff.Cap())

	var t []byte
	t = []byte("helllo world")
	x := t[:0]
	fmt.Println("len", len(t), "cap", cap(t))
	fmt.Println("len", len(x), "cap", cap(x))

	lis, err := kcp.ListenWithOptions(":10000", nil, 10, 3)

}
