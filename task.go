package httpserver

// TaskManager record request info, process into storage
// generally use for generate task id for asynchronous func
type TaskManager interface {
	Start(string) (int, error) // start task - return taskID
	SetState(int, int, string) error
}

type asyncTask struct {
	taskID int
	byteIn []byte
	doFunc func([]byte) error
}
