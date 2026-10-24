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

var defaultTransport = &http.Transport{
	Proxy:               http.ProxyFromEnvironment,
	MaxIdleConns:        10000,
	MaxIdleConnsPerHost: 1000,
	IdleConnTimeout:     90 * time.Second,
}

var defaultClient = &http.Client{
	Transport: defaultTransport,
	Timeout:   30 * time.Second,
}

// 32KB buffer that can be shared
var bufferPool = sync.Pool{
	New: func() any {
		b := make([]byte, 32*1024)
		return &b
	},
}

// Workers are taking Requests till the Channel is open.
type Worker struct {
	Id string
}

func (w *Worker) Start(ProxySettings *configParser.Proxy, RequestsChan <-chan request.RequestsJob, wg *sync.WaitGroup) {
	defer wg.Done()

	var jobWg sync.WaitGroup

	for Req := range RequestsChan {
		// made a copy
		job := Req
		jobWg.Add(1)
		go func() {
			defer jobWg.Done()
			w.handleJob(&job, defaultClient, ProxySettings)
		}()
	}

	jobWg.Wait()
	fmt.Println("Worker stopped")
}

func (w *Worker) handleJob(job *request.RequestsJob, client *http.Client, ProxySettings *configParser.Proxy) {
	defer close(job.Done)

	// Get the path to the response is being sent to, .
	// The algorithm is based on the configuration in the config.
	// The work of load balancing management is also managed by the router.
	forwardPath := router.GetPath(job, ProxySettings)
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

	buf := bufferPool.Get().(*[]byte)
	defer bufferPool.Put(buf)
	io.CopyBuffer(job.Writer, resp.Body, *buf)
}

func copyHeaders(src, dest http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dest.Add(k, v)
		}
	}
}
