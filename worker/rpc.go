package worker

import (
	"bytes"
	//"fmt"
	"io"
	"io/ioutil"
	"net/http"
	//"os"
	//"os/exec"
	//"path"
	//"strconv"
	//"strings"
	"time"

	//redis "gopkg.in/redis.v3"

	//"github.com/flike/golog"
	//"github.com/huyueMiss/kingtask/config"
	"github.com/huyueMiss/kingtask/core/errors"
	"github.com/huyueMiss/kingtask/task"



)

func (w *Worker) DoRpcTaskRequest(req *task.TaskRequest) (string, error) {
	var method string
	switch req.TaskType {
	case task.RpcTaskGET:
		method = "GET"
	case task.RpcTaskPOST:
		method = "POST"
	case task.RpcTaskPUT:
		method = "PUT"
	case task.RpcTaskDELETE:
		method = "DELETE"
	default:
		method = "GET"
	}
	url := req.BinName
	args := req.Args
	request, err := w.newHttpRequest(method, url, args)
	if err != nil {
		return "", err
	}
	result, err := w.callRpc(request, time.Second*time.Duration(req.MaxRunTime))
	return result, err
}

func (w *Worker) callRpc(req *http.Request, maxRunTime time.Duration) (string, error) {
	var timeout time.Duration
	if w.cfg.TaskRunTime != 0 {
		timeout = time.Duration(w.cfg.TaskRunTime) * time.Second
	} else {
		timeout = maxRunTime
	}

	//new a http client with timeout setting
	client := &http.Client{
		Timeout: timeout,
	}

	r, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	if r.StatusCode != http.StatusOK {
		return "", errors.NewError(string(buf))
	}

	return string(buf), nil
}

func (w *Worker) newHttpRequest(method string, url string, args string) (*http.Request, error) {
	var body io.Reader
	if len(args) != 0 {
		body = bytes.NewBuffer([]byte(args))
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}