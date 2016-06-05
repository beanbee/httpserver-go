package httpserver

import (
	"time"

	"testing"
)

var threadNum = 20

func TestHttpServer(t *testing.T) {
	server := NewServer("dbpserver", 5858).SetAsyncNum(threadNum)

	// handler sync http request
	server.HandlerRequst("POST", "/sync", syncDemo)

	// handler async http request
	server.HandlerAsyncRequst("POST", "/async", asyncDemo)
	// go func() {
	// 	if err := server.Start(); err != nil {
	// 		t.Fatalf("server start failed: %v", err)
	// 	}
	// }()

	// defer server.Stop()
}

func syncDemo(jsonIn []byte) (jsonOut []byte, err error) {
	return jsonIn, nil
}

func asyncDemo(jsonIn []byte) (err error) {
	time.Sleep(5 * time.Second)
	return nil
}
