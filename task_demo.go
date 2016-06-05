package httpserver

const (
	DEFAULT_TASK_ID     = -1
	TASK_STATE_RECEIVED = iota
	TASK_STATE_RUNNING
	TASK_STATE_SUCCESS
	TASK_STATE_FAILED
)

// default task manager
type defaultTaskManager struct{}

// record request url, return task id
func (d *defaultTaskManager) Start(reqUrl string) (int, error) {
	return DEFAULT_TASK_ID, nil
}

// update task state
func (d *defaultTaskManager) SetState(id, state int, result string) error {
	return nil
}
