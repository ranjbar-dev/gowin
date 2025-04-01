package data

import (
	"sync"

	"github.com/ranjbar-dev/gowin/types"
)

type Data struct {
	// clients connected to server
	clients      map[string]int64
	clientsMutex sync.Mutex
	// jobs to do
	jobs      map[string][]types.Job
	jobsMutex sync.Mutex
	// logs
	logs      []string
	logsMutex sync.Mutex
}

func NewData() *Data {

	return &Data{
		clients:      make(map[string]int64),
		clientsMutex: sync.Mutex{},
		jobs:         make(map[string][]types.Job),
		jobsMutex:    sync.Mutex{},
		logs:         []string{},
		logsMutex:    sync.Mutex{},
	}
}
