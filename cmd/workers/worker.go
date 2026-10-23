package workers

import (
	request "JerryCop/cmd/request"
	configParser "JerryCop/internal/configParser"
	router "JerryCop/internal/router"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// Workers are taking Requests till the Channel is open.
type Worker struct {
	Id string
}

func (w *Worker) Start(LocalProxy *configParser.Proxy, RequestsChan <-chan request.RequestsJob, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	for Req := range RequestsChan {
		w.handleJob(&Req, client, LocalProxy)
	}

	fmt.Println("Worker stopped")
}

func (w *Worker) handleJob(job *request.RequestsJob, client *http.Client, localProxy *configParser.Proxy) {
	defer close(job.Done)

	forwardPath := router.GetPath(job, localProxy)
	forwardRequest(job, client, forwardPath)
}

func forwardRequest(job *request.RequestsJob, client *http.Client, targetURL string) {
	upstreamReq, err := http.NewRequest(job.Req.Method, targetURL, job.Req.Body)
	if err != nil {
		errors.New("failed to ")
		return
	}

	copyHeaders(job.Req.Header, upstreamReq.Header)

	resp, err := client.Do(upstreamReq)
	if err != nil {
		errors.New("failed to execute")
		return
	}
	defer resp.Body.Close()

	copyHeaders(resp.Header, job.Writer.Header())
	job.Writer.WriteHeader(resp.StatusCode)
	io.Copy(job.Writer, resp.Body)
}

func copyHeaders(src, dst http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
