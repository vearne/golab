package main
import (
	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
	"crypto/sha1"
	"net"
	"time"
	"os"
	"os/signal"
	"syscall"
	"log"
)

const portEcho = "127.0.0.1:9999"
const portSink = "127.0.0.1:19999"
var key = []byte("testkey")
var pass = pbkdf2.Key(key, []byte(portSink), 4096, 32, sha1.New)

func listenEcho() (net.Listener, error) {
	//block, _ := NewNoneBlockCrypt(pass)
	//block, _ := NewSimpleXORBlockCrypt(pass)
	//block, _ := NewTEABlockCrypt(pass[:16])
	//block, _ := NewAESBlockCrypt(pass)
	//block, _ := kcp.NewSalsa20BlockCrypt(pass)
	return kcp.ListenWithOptions(portEcho, nil, 10, 3)
}

func handleEcho(conn *kcp.UDPSession) {
	conn.SetStreamMode(true)
	conn.SetWindowSize(4096, 4096)
	conn.SetNoDelay(1, 10, 2, 1)
	conn.SetDSCP(46)
	conn.SetMtu(1400)
	conn.SetACKNoDelay(false)
	conn.SetReadDeadline(time.Now().Add(time.Minute))
	conn.SetWriteDeadline(time.Now().Add(time.Minute))
	buf := make([]byte, 65536)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}
		conn.Write(buf[:n])
	}
}

func echoServer() {
	l, err := listenEcho()
	if err != nil {
		panic(err)
	}

	go func() {
		kcplistener := l.(*kcp.Listener)
		kcplistener.SetReadBuffer(4 * 1024 * 1024)
		kcplistener.SetWriteBuffer(4 * 1024 * 1024)
		kcplistener.SetDSCP(46)
		for {
			s, err := l.Accept()
			if err != nil {
				return
			}

			// coverage test
			s.(*kcp.UDPSession).SetReadBuffer(4 * 1024 * 1024)
			s.(*kcp.UDPSession).SetWriteBuffer(4 * 1024 * 1024)
			go handleEcho(s.(*kcp.UDPSession))
		}
	}()
}

func main() {
	echoServer()
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	sig := <-ch
	log.Println("get signal", sig)
}