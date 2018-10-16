package main
import (
	"github.com/xtaci/kcp-go"
	"fmt"
	"time"
	"golang.org/x/crypto/pbkdf2"
	"crypto/sha1"
)

const portEcho = "127.0.0.1:9999"
const portSink = "127.0.0.1:19999"
var key = []byte("testkey")
var pass = pbkdf2.Key(key, []byte(portSink), 4096, 32, sha1.New)

func dialEcho() (*kcp.UDPSession, error) {
	//block, _ := NewNoneBlockCrypt(pass)
	//block, _ := NewSimpleXORBlockCrypt(pass)
	//block, _ := NewTEABlockCrypt(pass[:16])
	//block, _ := NewAESBlockCrypt(pass)
	//block, _ := kcp.NewSalsa20BlockCrypt(pass)
	sess, err := kcp.DialWithOptions(portEcho, nil, 10, 3)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	sess.SetStreamMode(true)
	sess.SetStreamMode(false)
	sess.SetStreamMode(true)
	sess.SetWindowSize(4096, 4096)
	sess.SetReadBuffer(4 * 1024 * 1024)
	sess.SetWriteBuffer(4 * 1024 * 1024)
	sess.SetStreamMode(true)
	sess.SetNoDelay(1, 10, 2, 1)
	sess.SetMtu(1400)
	sess.SetMtu(1600)
	sess.SetMtu(1400)
	sess.SetACKNoDelay(true)
	sess.SetDeadline(time.Now().Add(time.Minute))
	return sess, err
}


func main(){
	cli, err := dialEcho()
	if err != nil {
		panic(err)
	}
	cli.SetWriteDelay(true)
	cli.SetDUP(1)
	const N = 100
	buf := make([]byte, 10)
	for i := 0; i < N; i++ {
		msg := fmt.Sprintf("hello%v", i)
		cli.Write([]byte(msg))
		if n, err := cli.Read(buf); err == nil {
			fmt.Println(string(buf[:n]))
			if string(buf[:n]) != msg {
				fmt.Println("不一致", string(buf[:n]),  msg)
			}
		} else {
			panic(err)
		}

	}
	cli.Close()
}