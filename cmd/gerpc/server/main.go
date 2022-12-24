package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"geek/gerpc"
	"geek/gerpc/codec"
	"geek/glog"
)

func startServer(addr chan string) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		glog.Errorf("server listen error: %v", err)
	}
	glog.Infof("server listening on %v", l.Addr())
	addr <- l.Addr().String()
	gerpc.Accept(l)
}

func main() {
	addr := make(chan string)
	go startServer(addr)

	conn, err := net.Dial("tcp", <-addr)
	if err != nil {
		glog.Errorf("dial connection error: %v", err)
		return
	}
	defer conn.Close()

	time.Sleep(time.Second)
	json.NewEncoder(conn).Encode(gerpc.DefaultOption)
	codecer := codec.NewGobCodec(conn)
	for i := 0; i < 5; i++ {
		header := &codec.Header{
			ServiceMethod: "Index.Sum",
			Seq:           uint64(i),
		}
		codecer.Write(header, fmt.Sprintf("gerpc req %d", header.Seq))
		codecer.ReadHeader(header)

		// reply := ""
		var reply string
		if err := codecer.ReadBody(&reply); err != nil {
			glog.Errorf("read body error: %v", err)
		}
		glog.Infof("reply: %s", reply)
	}
}
