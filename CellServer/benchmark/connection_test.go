package benchmark

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"testing"
)

var addr = flag.String("addr", "127.0.0.1:36000", "http service address")

func BenchmarkConnection(b *testing.B) {

	flag.Parse()
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Fatal("dial:", err)
		}
		c.Close()
	}
}
