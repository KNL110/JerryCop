package protocol

import (
	"net"
	"sync"
	"errors"
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

func (t *TCPListener) AcceptLoop (LocalMaxConcListener *sync.WaitGroup) Error{
	defer LocalMaxConcListener.Done()
	ln, err := net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return errors.New("Encountered an error while listening ")
	}
	for {
		
	}
}

func (*TCPListener) Listen (MaxConcReq int, LocalMaxConcListener *sync.WaitGroup) {
	if (int(LocalMaxConcListener.state[0:32]) <= MaxConcReq) {
		LocalMaxConcListener.Add(1)
		go actuallyListen()

		
}