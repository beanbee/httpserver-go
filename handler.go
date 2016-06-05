package httpserver

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// request callback func
type reqCall struct {
	Method   string
	HttpFunc func(rw http.ResponseWriter, req *http.Request, byteIn []byte)
}

// http handler
type defaultHandler struct {
	Mux map[string]reqCall // method_name - request[post] function
}

// http server handler - implement
func (def *defaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := def.Mux[r.RequestURI]; ok && h.Method == r.Method {
		byteIn, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		h.HttpFunc(w, r, byteIn)
	}
}

// basic return json
// format: {"code": xxx, "msg": "", "id": 0}
func retJson(rw http.ResponseWriter, code, taskid int, msg string) {
	content := fmt.Sprintf(`{"code": %d, "msg": "%s"%s}`, code, msg,
		func(id int) string {
			if id != DEFAULT_TASK_ID {
				return fmt.Sprintf(`, "id": %d`, id)
			}
			return ""
		}(taskid))
	http.Error(rw, content, code)
}
