package main

import (
	request "JerryCop/cmd/request"
	workers "JerryCop/cmd/workers"
	configparser "JerryCop/internal/configParser"
	"JerryCop/internal/protocol"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/google/uuid"

	_ "github.com/joho/godotenv"
)

func main() {

	// Step1: Load Config
	LocalProxy, err := configparser.LoadConfig()

	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	// Step2: Create Dedicated Channel
	RequestsChan := make(chan request.RequestsJob, LocalProxy.MaxConcReq)
	// RequestsChan := make(chan []byte, 10000)

	// Step3: Spawn Workers
	jobWorkers := make([]*workers.Worker, LocalProxy.NumWorkers)
	var workersWg sync.WaitGroup

	for i := range jobWorkers {
		jobWorkers[i] = &(workers.Worker{Id: uuid.NewString()})
		workersWg.Add(1)
		go jobWorkers[i].Start(&LocalProxy, RequestsChan, &workersWg)
	}

	// Step4: spawns Listener and listens to the incoming requests.
	Listener := protocol.NewHTTPListener(LocalProxy.ProxyPort, LocalProxy.MaxConcReq, RequestsChan)
	if Listener == nil {
		fmt.Printf("Error creating Listener for port %s\n", LocalProxy.ProxyPort)
		os.Exit(1)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := Listener.Listen(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Listener error:", err)
		}
	}()

	fmt.Println("Blocking, press ctrl+c to continue...")
	<-done
	fmt.Println("Interrupted! Shutting down...")

	// Stop the Listener
	_ = Listener.Server.Close()

	// Tell the workers threads to stop taking the requests.
	close(RequestsChan)
	// Wait till the worker threads are done with their workload.
	workersWg.Wait()

	fmt.Println("All workers finished. Exiting.")
}
