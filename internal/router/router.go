package router

import (
	"JerryCop/internal/request"
	configparser "JerryCop/internal/configParser"
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

func GetPath(job *request.RequestsJob, proxySettings *configparser.Proxy, state ...*State) (string, error) {
	var s *State
	if len(state) > 0 {
		s = state[0]
	}

	switch proxySettings.DistAlgo {
	case "RANDOM":
		return getPathRandom(job, proxySettings, s)
	}

	return "", errors.New("could not get path.")
}

