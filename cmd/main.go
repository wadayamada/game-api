package main

import (
	"flag"
	"log"

	"20dojo-online/pkg/server"
)

var (
	// Listenするアドレス+ポート
	addr string
)

func init() {
	flag.StringVar(&addr, "addr", ":8080", "tcp host:port to connect")
	flag.Parse()
}

func main() {
	//ファイル名もログに残す設定
	log.SetFlags(log.LstdFlags | log.Llongfile)
	server.Serve(addr)
}
