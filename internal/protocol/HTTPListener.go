package protocol

import (
	"JerryCop/internal/request"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
)

var reqCounter atomic.Uint64

type HTTPListener struct {
	Server       *http.Server
	RequestsChan chan<- request.RequestsJob
}

func NewHTTPListener(listenAddr string, MaxConcReq int, RequestsChan chan<- request.RequestsJob) *HTTPListener {
	if RequestsChan == nil {
		return nil
	}

	mux := mux.NewRouter()

	listener := &HTTPListener{
		RequestsChan: RequestsChan,
	}

	mux.PathPrefix("/").HandlerFunc(listener.handleRequest)

	Limit := limitMiddleware(MaxConcReq, mux)
	listener.Server = &http.Server{
		Addr:         listenAddr,
		Handler:      Limit,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	return listener
}

// func (s *HTTPListener) handleRequest(w http.ResponseWriter, r *http.Request) {
// 	done := make(chan struct{})
// 	job := request.RequestsJob{
// 		ID:     strconv.FormatUint(reqCounter.Add(1), 10),
// 		Req:    r,
// 		Writer: w,
// 		Done:   done,
// 	}

// 	s.RequestsChan <- job
// 	<-done
// }

var doneChanPool = sync.Pool{
    New: func() any {
        return make(chan struct{})
    },
}

func (s *HTTPListener) handleRequest(w http.ResponseWriter, r *http.Request) {
    done := doneChanPool.Get().(chan struct{})
    
    job := request.RequestsJob{
        ID:     strconv.FormatUint(reqCounter.Add(1), 10),
        Req:    r,
        Writer: w,
        Done:   done,
    }

    s.RequestsChan <- job
    <-done

    // Put back into pool for reuse (must be drained or non-closed)
    // Note: since handleJob closes job.Done, you recreate or use a 1-capacity buffered struct{}:
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

func (s *HTTPListener) Listen() error {

	if err := s.Server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
