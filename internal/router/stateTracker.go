package router

import (
	configparser "JerryCop/internal/configParser"
	"sync/atomic"
)

// Although the config and stateTracker both deal with the servers,
// config only loads the data about the servers in memory
// While the stateTracker deals with the information
// tracking about the Connections and Servers at the runtime.

type State struct {
	LastConnection   InternalServer
	IsInternalServer map[string]struct{}
	NumOfConns       map[InternalServer]int
	TotalServers     int
	RoundRobinIndex  map[string]*atomic.Uint64
}

type InternalServer struct {
	ServerPath string
}

func (State) NewState(configSettings *configparser.Proxy) State {
	state := State{
		LastConnection:   InternalServer{},
		IsInternalServer: make(map[string]struct{}),
		NumOfConns:       make(map[InternalServer]int),
		TotalServers:     len(configSettings.ServerConfig),
		RoundRobinIndex:  make(map[string]*atomic.Uint64),
	}

	for _, server := range configSettings.ServerConfig {
		prefix := server.ServerRouting.Prefix
		state.RoundRobinIndex[prefix] = new(atomic.Uint64)
	}

	return state
}
