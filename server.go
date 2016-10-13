package httpserver

import (
	"fmt"
	"net/http"
	"time"
)

const (
	DEFAULT_ASYNC_NUM    = 4
	DEFAULT_STOP_TIMEOUT = 1800 * time.Second
)

type Server struct {
	name       string
	port       int
	asyncNum   int                // max goroutine number for async tasks
	handlerMap map[string]reqCall // map[path]reqCall
	taskPipe   chan asyncTask     // channel for async tasks
	donePipe   chan interface{}   // wait all async tasks
	endSignal  chan interface{}   // perform gracefully stop

	listener    *keepAliveListener
	taskManager TaskManager
}

// create new Server using default settings
func NewServer(name string, port int) *Server {
	return &Server{
		name:       name,
		port:       port,
		asyncNum:   DEFAULT_ASYNC_NUM,
		handlerMap: make(map[string]reqCall),
		taskPipe:   make(chan asyncTask, DEFAULT_ASYNC_NUM),
		donePipe:   make(chan interface{}, DEFAULT_ASYNC_NUM),
		endSignal:  make(chan interface{}),

		taskManager: new(defaultTaskManager),
	}
}

// set max async thread number
func (s *Server) SetAsyncNum(num int) *Server {
	s.asyncNum = num
	s.taskPipe = make(chan asyncTask, num)
	s.donePipe = make(chan interface{}, num)
	return s
}

// set task manager
func (s *Server) SetTaskManager(taskManager TaskManager) *Server {
	s.taskManager = taskManager
	return s
}

// handler sync http request
// input:  method, url string, syncFunc func([]byte) ([]byte, error)
// ourput: nil
// NOTE:   sync request will not be recorded in taskManager
func (s *Server) HandlerRequst(method, url string, syncFunc func([]byte) ([]byte, error)) {
	httpFunc := func(rw http.ResponseWriter, req *http.Request, byteIn []byte) {
		byteOut, err := syncFunc(byteIn)
		if err != nil {
			retJson(rw, http.StatusInternalServerError, DEFAULT_TASK_ID, err.Error())
			return
		}
		http.Error(rw, string(byteOut), http.StatusOK)
	}

	// add to http handler map
	s.handlerMap[url] = reqCall{
		Method:   method,
		HttpFunc: httpFunc,
	}
}

// handler async http request
// input:  method, url string, syncFunc func([]byte) ([]byte, error)
// ourput: nil
// NOTE:   create task using taskManager, pass taskID to async task channel through struct
func (s *Server) HandlerAsyncRequst(method, url string, asyncFunc func([]byte) error) {
	httpFunc := func(rw http.ResponseWriter, req *http.Request, byteIn []byte) {
		taskID, err := s.taskManager.Start(url)
		if err != nil {
			retJson(rw, http.StatusInternalServerError, DEFAULT_TASK_ID, fmt.Sprintf("create task failed: %v", err))
			return
		}

		// add to task pipe
		s.taskPipe <- asyncTask{
			taskID: taskID,
			byteIn: byteIn,
			doFunc: asyncFunc,
		}
		retJson(rw, http.StatusOK, taskID, "start task success")
	}

	// add to http handler map
	s.handlerMap[url] = reqCall{
		Method:   method,
		HttpFunc: httpFunc,
	}
}

// start http server
func (s *Server) Start() (err error) {
	server := http.Server{
		Handler: &defaultHandler{
			Mux: s.handlerMap,
		},
		ReadTimeout: 30 * time.Second, // to prevent abuse of "keep-alive" requests by clients
	}

	// create tcp listener
	s.listener, err = newKeepAliveListener(s.port)
	if err != nil {
		return err
	}

	// perform async task with background goroutines
	for i := 0; i < s.asyncNum; i++ {
		go func() {
			for task := range s.taskPipe {
				myfunc := task.doFunc
				s.taskManager.SetState(task.taskID, TASK_STATE_RUNNING, "task running")
				if err := myfunc(task.byteIn); err != nil {
					s.taskManager.SetState(task.taskID, TASK_STATE_FAILED, err.Error())
				} else {
					s.taskManager.SetState(task.taskID, TASK_STATE_SUCCESS, "task end successfully")
				}
			}
			s.donePipe <- struct{}{}
		}()
	}

	// collect endsignal and close task pipe
	go func() {
		<-s.endSignal
		s.listener.Close() // close listener
		timeOut := time.Now().Add(DEFAULT_STOP_TIMEOUT)
		for i := 0; time.Now().Before(timeOut); i++ {
			if len(s.taskPipe) == 0 {
				break
			}
		}
		close(s.taskPipe)
	}()

	return server.Serve(s.listener)
}

// gracefully stop server
func (s *Server) Stop() {
	close(s.endSignal)
	if s.listener != nil {
		// awaitCompletion
		for i := 0; i < s.asyncNum; i++ {
			<-s.donePipe
			fmt.Printf("[%d] async goroutine closed\n", i+1)
		}
	}
	close(s.donePipe)
}
