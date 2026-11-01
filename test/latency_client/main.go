package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"sort"
	"sync"
	"time"
)

type LatencyResult struct {
	Duration time.Duration
	Status   int
	Err      error
}

func main() {
	targetURL := flag.String("url", "http://localhost:8080/api/users?id=123", "target URL to benchmark")
	totalRequests := flag.Int("n", 100, "total number of requests")
	concurrency := flag.Int("c", 10, "number of concurrent workers")
	flag.Parse()

	if *totalRequests <= 0 || *concurrency <= 0 {
		fmt.Println("Total requests (-n) and concurrency (-c) must be > 0")
		return
	}

	fmt.Printf("=== JerryCop Latency & Benchmark Tool ===\n")
	fmt.Printf("Target URL:    %s\n", *targetURL)
	fmt.Printf("Requests:      %d\n", *totalRequests)
	fmt.Printf("Concurrency:   %d\n\n", *concurrency)

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        1000,
			MaxIdleConnsPerHost: 1000,
			IdleConnTimeout:     30 * time.Second,
		},
	}

	// Warmup / sanity probe
	probeStart := time.Now()
	probeResp, err := client.Get(*targetURL)
	if err != nil {
		fmt.Printf("[PROBE FAILED] Cannot reach %s: %v\n", *targetURL, err)
		fmt.Println("Make sure JerryCop proxy and backend servers are running!")
		return
	}
	body, _ := io.ReadAll(probeResp.Body)
	probeResp.Body.Close()
	fmt.Printf("[PROBE OK] Status: %d, Round-trip: %v\n", probeResp.StatusCode, time.Since(probeStart))
	fmt.Printf("[PROBE BODY] %s\n\n", string(body))

	results := make([]LatencyResult, *totalRequests)
	jobs := make(chan int, *totalRequests)
	for i := 0; i < *totalRequests; i++ {
		jobs <- i
	}
	close(jobs)

	var wg sync.WaitGroup
	startTime := time.Now()

	for w := 0; w < *concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				t0 := time.Now()
				resp, err := client.Get(*targetURL)
				elapsed := time.Since(t0)

				res := LatencyResult{Duration: elapsed, Err: err}
				if err == nil {
					res.Status = resp.StatusCode
					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()
				}
				results[idx] = res
			}
		}()
	}

	wg.Wait()
	totalDuration := time.Since(startTime)

	// Process stats
	var durations []time.Duration
	statusCounts := make(map[int]int)
	errCount := 0
	var totalElapsed time.Duration

	for _, r := range results {
		if r.Err != nil {
			errCount++
			continue
		}
		statusCounts[r.Status]++
		durations = append(durations, r.Duration)
		totalElapsed += r.Duration
	}

	if len(durations) == 0 {
		fmt.Printf("All %d requests failed!\n", *totalRequests)
		return
	}

	sort.Slice(durations, func(i, j int) bool {
		return durations[i] < durations[j]
	})

	rps := float64(len(durations)) / totalDuration.Seconds()
	avg := totalElapsed / time.Duration(len(durations))
	min := durations[0]
	max := durations[len(durations)-1]
	p50 := durations[int(float64(len(durations))*0.50)]
	p90 := durations[int(float64(len(durations))*0.90)]
	p95 := durations[int(float64(len(durations))*0.95)]
	p99 := durations[int(float64(len(durations))*0.99)]

	fmt.Println("--- Benchmark Summary ---")
	fmt.Printf("Completed Requests: %d / %d (Errors: %d)\n", len(durations), *totalRequests, errCount)
	fmt.Printf("Total Test Time:    %v\n", totalDuration)
	fmt.Printf("Throughput:         %.2f req/sec\n\n", rps)

	fmt.Println("--- Latency Breakdown ---")
	fmt.Printf("Min:     %v\n", min)
	fmt.Printf("Average: %v\n", avg)
	fmt.Printf("Median:  %v\n", p50)
	fmt.Printf("P90:     %v\n", p90)
	fmt.Printf("P95:     %v\n", p95)
	fmt.Printf("P99:     %v\n", p99)
	fmt.Printf("Max:     %v\n\n", max)

	fmt.Println("--- Status Code Distribution ---")
	for status, count := range statusCounts {
		fmt.Printf("HTTP %d: %d\n", status, count)
	}
}
