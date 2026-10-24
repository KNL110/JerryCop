package router

import (
	configparser "JerryCop/internal/configParser"
)

// Although the config and stateTracker both deal with the servers,
// config only loads the data about the servers in memory
// While the stateTracker deals with the information
// tracking about the Connections and Servers at the runtime.

type State struct {
	LastConnection   InternalServer
	IsInternalServer map[string]bool
	NumOfConns       map[InternalServer]int
	TotalServers     int
}

type InternalServer struct {
	ServerPath string
}

func (State) NewState(configSettings *configparser.Proxy) State {
	var NumofConns map[InternalServer]int
	for server := range 
	
	state := State{
		
	}

	return state
}
