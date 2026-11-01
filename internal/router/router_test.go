package router

import (
	"JerryCop/internal/request"
	configparser "JerryCop/internal/configParser"
	"net/http/httptest"
	"testing"
)

func TestGetPath_Random(t *testing.T) {
	cfg := &configparser.Proxy{
		DistAlgo: "RANDOM",
		ServerConfig: []configparser.ServerConfig{
			{
				ServerRouting: configparser.Routing{
					Prefix:    "/api",
					ProxyPass: []string{"http://10.0.0.1:8000", "http://10.0.0.2:8000"},
				},
			},
		},
	}

	req := httptest.NewRequest("GET", "/api/v1/users?active=true", nil)
	job := &request.RequestsJob{Req: req}

	target, err := GetPath(job, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedPrefix1 := "http://10.0.0.1:8000/api/v1/users?active=true"
	expectedPrefix2 := "http://10.0.0.2:8000/api/v1/users?active=true"
	if target != expectedPrefix1 && target != expectedPrefix2 {
		t.Errorf("got %s, expected either %s or %s", target, expectedPrefix1, expectedPrefix2)
	}
}

func TestGetPath_NoMatch(t *testing.T) {
	cfg := &configparser.Proxy{
		DistAlgo: "RANDOM",
		ServerConfig: []configparser.ServerConfig{
			{
				ServerRouting: configparser.Routing{
					Prefix:    "/admin",
					ProxyPass: []string{"http://10.0.0.1:8000"},
				},
			},
		},
	}

	req := httptest.NewRequest("GET", "/public/index.html", nil)
	job := &request.RequestsJob{Req: req}

	_, err := GetPath(job, cfg)
	if err == nil {
		t.Fatal("expected error for unmatched route, got nil")
	}
}

func TestGetPath_UnsupportedAlgo(t *testing.T) {
	cfg := &configparser.Proxy{
		DistAlgo: "LEAST_CONN",
		ServerConfig: []configparser.ServerConfig{
			{
				ServerRouting: configparser.Routing{
					Prefix:    "/",
					ProxyPass: []string{"http://10.0.0.1:8000"},
				},
			},
		},
	}

	req := httptest.NewRequest("GET", "/test", nil)
	job := &request.RequestsJob{Req: req}

	_, err := GetPath(job, cfg)
	if err == nil {
		t.Fatal("expected error for unsupported algorithm, got nil")
	}
}

func TestGetPath_RoundRobin(t *testing.T) {
	cfg := &configparser.Proxy{
		DistAlgo: "ROUND_ROBIN",
		ServerConfig: []configparser.ServerConfig{
			{
				ServerRouting: configparser.Routing{
					Prefix:    "/api",
					ProxyPass: []string{"http://10.0.0.1:8000", "http://10.0.0.2:8000"},
				},
			},
		},
	}

	state := State{}.NewState(cfg)

	req1 := httptest.NewRequest("GET", "/api/data", nil)
	target1, err := GetPath(&request.RequestsJob{Req: req1}, cfg, &state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req2 := httptest.NewRequest("GET", "/api/data", nil)
	target2, err := GetPath(&request.RequestsJob{Req: req2}, cfg, &state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if target1 != "http://10.0.0.1:8000/api/data" {
		t.Errorf("expected target1 to be 10.0.0.1, got %s", target1)
	}
	if target2 != "http://10.0.0.2:8000/api/data" {
		t.Errorf("expected target2 to be 10.0.0.2, got %s", target2)
	}
}
