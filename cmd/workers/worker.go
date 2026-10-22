package workers

import (
	configparser "JerryCop/internal/configParser"
)

// these workers are supposed to be open indefinitely as long as channel is open

type Worker struct {
	Id string
}

func (w *Worker) Start(LocalProxy *configparser.Proxy, RequestsChan <-chan []byte, wg *sync.WaitGroup) {
	for Req := range RequestsChan {
		// actual internal working
	}
	wg.Done()
}
