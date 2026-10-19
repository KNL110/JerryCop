package protocol

import (
	"net"
	"sync"
)

type TCPListener struct {
	ListenAddr string
	Listener net.Listener
	RequestsChan chan<- []byte
}

func NewTCPListener (listenaddress ListenAddr) TCPListener {
	NewTCPListener := TCPListener{
		ListenAddr: listenaddress,
	}
}

func (*TCPListener) actuallyListen (LocalMaxConcListener *sync.WaitGroup) {
	defer LocalMaxConcListener.Done()

	// listen to requests here.
	...
}

func (*TCPListener) Listen (MaxConcReq int, LocalMaxConcListener *sync.WaitGroup) {
	if (int(LocalMaxConcListener.state[0:32]) <= MaxConcReq) {
		LocalMaxConcListener.Add(1)
		go actuallyListen()
	}
}