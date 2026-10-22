package protocol

import (
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"JerryCop/cmd/request"
)


type HTTPListener struct {
	Server *http.Server
	RequestsChan chan<- request.RequestsJob
}

func NewHTTPListener(listenAddr string, MaxConcReq int, RequestsChan chan<- request.RequestsJob) *HTTPListener {
	mux := mux.NewRouter()

	if RequestsChan == nil {
		return nil
	}

	Limit := limitMiddleware(MaxConcReq,mux)
	NewHTTPListener := HTTPListener{
		Server : &http.Server{
			Addr: listenAddr,
			Handler: Limit,
			WriteTimeout: time.Second * 30,
			ReadTimeout:  time.Second * 10,
			IdleTimeout:  time.Minute,
		},
		RequestsChan: RequestsChan,
	}

	return &NewHTTPListener
}

func limitMiddleware(maxConcurrent int, next http.Handler) http.Handler {
	// A buffered channel acts as an semaphore
	semaphore := make(chan struct{}, maxConcurrent)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Acquire a slot. If full, this blocks until a slot opens up.
		semaphore <- struct{}{}
		
		// Ensure the slot is released when the request finishes
		defer func() { <-semaphore }()

		next.ServeHTTP(w, r)
	})
}

func (s *HTTPListener)Listen() error{

	if err := s.Server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
