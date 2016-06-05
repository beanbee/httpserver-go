# httpserver-go
restful http server with async-request handler &amp; stoppable listener (Golang)

Usage:

import (
        httpserver "github.com/beanbee/httpserver-go"
)

func main() {
<!----> create new http server
        server := httpserver.NewServer("server", 3005).SetAsyncNum(20)
        
	// handler sync http request
	server.HandlerRequst("POST", "/sync", syncDemo)

	// handler async http request
	server.HandlerAsyncRequst("POST", "/async", asyncDemo)

         go func() {
        if err := server.Start(); err != nil {
                log.Printf("server failed: %v", err)
        }
 }()

        // stop agent with system signal
        EndChannel := make(chan os.Signal)
        signal.Notify(EndChannel, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
        select {
        case output := <-EndChannel:
                log.Printf("end myserver process: %s", output)
                server.Stop()
                break
        }
        close(EndChannel)
        log.Printf("all work done")
//      time.Sleep(20 * time.Second)
}

func syncDemo(jsonIn []byte) (jsonOut []byte, err error) {
        log.Printf("[syncDemo] jsonIn: %v", string(jsonIn[:]))

        return jsonIn, nil
}

func asyncDemo(jsonIn []byte) (err error) {
        time.Sleep(5 * time.Second)
        log.Printf("[asyncDemo] jsonIn: %v", string(jsonIn[:]))

        return nil
}
