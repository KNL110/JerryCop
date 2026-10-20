package workers

import (
	configparser "JerryCop/internal/configParser"
)

// these workers are supposed to be open indefinitely as long as channel is open

type Worker struct {
	Id string
}

func (*Worker) Start(LocalProxy *configparser.Proxy, RequestsChan <-chan []byte) {
	for Req := range RequestsChan {
		// actual internal working
	}
}
