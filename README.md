# httpserver-go
restful http server with async-request handler &amp; stoppable listener (Golang)

Usage:

`
import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpserver "github.com/beanbee/httpserver-go"
)

func main() {
	// create new http server with max async 20 goroutines
	server := httpserver.NewServer("mytest", 3005).SetAsyncNum(20)

	// handler sync http request
	server.HandlerRequst("POST", "/sync", syncDemo)

	// handler async http request
	server.HandlerAsyncRequst("POST", "/async", asyncDemo)

	if err := server.Start(); err != nil {
		log.Printf("server failed: %v", err)
	}

    // await completion for all request
    server.Stop()

}

// simple handler for sync request
// get response data immediately
func syncDemo(jsonIn []byte) (jsonOut []byte, err error) {
	log.Printf("[syncDemo] jsonIn: %v", string(jsonIn[:]))

	return jsonIn, nil
}

// simple handler for async request
// return task info in response data when performing request handler asynchronously
func asyncDemo(jsonIn []byte) (err error) {
	time.Sleep(5 * time.Second)
	log.Printf("[asyncDemo] jsonIn: %v", string(jsonIn[:]))

	return nil
}
`
