package cmd

import (
	"JerryCop/cmd/workers"
	configparser "JerryCop/internal/configParser"
	"JerryCop/internal/protocol"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "github.com/joho/godotenv"
)

func main() {
	// Step1: Load Config
	LocalProxy := configparser.LoadConfig()

	// Step2: Create Dedicated Channel
	RequestsChan := make(chan []byte, 10000)

	// Step3: Spawn Workers
	workers := workers.Worker
	workersWg := sync.WaitGroup
	maxConcListener := sync.WaitGroup

	for _ := range LocalProxy.NumWorkers {
		go workers.StartWorker(&LocalProxy, RequestsChan)
		workersWg.add(1)
	}

	// Step4: spawn Listener
	Listener := protocol.TCPListener
	Listener.Listen()

	// Close the channel and exit the process when interrupted.
	// Will have to implement a more graceful shutdown later.
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	close(RequestsChan)
	fmt.Println("Blocking, press ctrl+c to continue...")
	<-done
	fmt.Println("\nInterrupted!")
}
