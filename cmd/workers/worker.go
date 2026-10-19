package workers

import (
	configparser "JerryCop/internal/configParser"
)

// these workers are supposed to be open indefinitely as long as channel is open

type Worker struct {
}

func (*Worker) StartWorker(LocalProxy *configparser.Proxy, RequestsChan <-chan []byte) {
	for Req := range RequestsChan {
		// actual internal working
	}
}
