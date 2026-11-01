package router

import (
	configparser "JerryCop/internal/configParser"
	"JerryCop/internal/request"
	"errors"
	"math/rand/v2"
)

func getPathRandom(job *request.RequestsJob, proxySettings *configparser.Proxy, state *State) (string, error) {
	server, found := FindMatchingServer(job.Req.URL.Path, proxySettings.ServerConfig)
	if !found || len(server.ServerRouting.ProxyPass) == 0 {
		return "", errors.New("no matching upstream server found for path")
	}

	randIndex := rand.IntN(len(server.ServerRouting.ProxyPass))
	chosenUpstream := server.ServerRouting.ProxyPass[randIndex]

	return BuildTargetURL(chosenUpstream, job.Req.URL.Path, job.Req.URL.RawQuery)
}

func getPathRoundRobin(job *request.RequestsJob, proxySettings *configparser.Proxy, state *State) (string, error) {
	server, found := FindMatchingServer(job.Req.URL.Path, proxySettings.ServerConfig)
	if !found || len(server.ServerRouting.ProxyPass) == 0 {
		return "", errors.New("no matching upstream server found for path")
	}

	upstreams := server.ServerRouting.ProxyPass
	prefix := server.ServerRouting.Prefix
	var chosenUpstream string
	if state != nil && state.RoundRobinIndex != nil && state.RoundRobinIndex[prefix] != nil {
		idx := state.RoundRobinIndex[prefix].Add(1) - 1
		chosenUpstream = upstreams[idx%uint64(len(upstreams))]
	} else {
		chosenUpstream = upstreams[0] // fallback
	}

	return BuildTargetURL(chosenUpstream, job.Req.URL.Path, job.Req.URL.RawQuery)
}

func GetPath(job *request.RequestsJob, proxySettings *configparser.Proxy, state ...*State) (string, error) {
	var s *State
	if len(state) > 0 {
		s = state[0]
	}

	switch proxySettings.DistAlgo {
	case "RANDOM":
		return getPathRandom(job, proxySettings, s)
	case "ROUND_ROBIN":
		return getPathRoundRobin(job, proxySettings, s)
	}

	return "", errors.New("could not get path.")
}
