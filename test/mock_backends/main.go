package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func startServer(port string, name string, delay time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if delay > 0 {
			time.Sleep(delay)
		}
		w.Header().Set("Content-Type", "application/json")
		resp := fmt.Sprintf(`{"server":"%s","port":"%s","path":"%s","query":"%s","time":"%s"}`,
			name, port, r.URL.Path, r.URL.RawQuery, time.Now().Format(time.RFC3339Nano))
		w.Write([]byte(resp))
	})

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("[Backend] %s listening on http://127.0.0.1:%s (simulated delay: %v)\n", name, port, delay)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("[%s] error: %v\n", name, err)
	}
}

func main() {
	delayMs := flag.Int("delay", 0, "simulated backend processing delay in milliseconds")
	flag.Parse()

	delay := time.Duration(*delayMs) * time.Millisecond
	var wg sync.WaitGroup

	backends := []struct {
		port string
		name string
	}{
		{"8001", "backend1"},
		{"8002", "backend2"},
		{"8003", "api_backend"},
	}

	fmt.Println("=== Starting Mock Internal Servers ===")
	for _, b := range backends {
		wg.Add(1)
		go startServer(b.port, b.name, delay, &wg)
	}

	wg.Wait()
}
