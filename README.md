# httpserver-go
restful http server with async-request handler &amp; stoppable listener (Golang)

Usage:

```Go
import (
	"log"
	"time"
	"os"
	"os/signal"
	"syscall"

	httpserver "github.com/beanbee/httpserver-go"
)

func main() {
	// create new http server with max async 20 goroutines
	server := httpserver.NewServer("mytest", 3005).SetAsyncNum(20)

	// handler sync http request
	server.HandlerRequst("POST", "/sync", syncDemo)

	// handler async http request
	server.HandlerAsyncRequst("POST", "/async", asyncDemo)
	
	go func(){
		if err := server.Start(); err != nil {
			log.Printf("server failed: %v", err)
	    	}
	}()

	// you can stop server using Stop() method which could await completion for all requests
	// finishing off some extra-works by a system signal is recommended
	EndChannel := make(chan os.Signal)
	signal.Notify(EndChannel, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	select {
	case output := <-EndChannel:
		log.Printf("end http server process by: %s", output)
		server.Stop()
	}
	close(EndChannel)
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
```
