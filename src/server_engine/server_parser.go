package server_engine

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//parse servers
type ParseServer interface {
	ParseServers() *ServerList
}

type ServerJson struct {
	Path       string
	ServerList ServerList
}

func (server *ServerJson) ParseServers() *ServerList {
	serverJson, err := os.Open(server.Path)
	if err != nil {
		panic(err)
	}
	defer serverJson.Close()

	//read
	byteJson, err := ioutil.ReadAll(serverJson)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal([]byte(byteJson), &server.ServerList); err != nil {
		panic(err)
	}

	return &server.ServerList
}
