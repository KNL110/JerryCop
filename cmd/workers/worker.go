package workers

import (
	request "JerryCop/cmd/request"
	configparser "JerryCop/internal/configParser"
	"fmt"
	"sync"
)

// these workers are supposed to be open indefinitely as long as channel is open

type Worker struct {
	Id string
}

func (w *Worker) Start(LocalProxy *configparser.Proxy, RequestsChan <-chan request.RequestsJob, wg *sync.WaitGroup) {
	for Req := range RequestsChan {
		// actual internal working
	}
	wg.Done()

	fmt.Println("Worker stopped")
}
