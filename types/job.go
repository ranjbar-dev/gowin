package types

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Job struct {
	ID        string   `json:"id"`
	Name      JobName  `json:"name"`
	Params    []string `json:"params"`
	Timestamp int64    `json:"timestamp"`
}

func NewJob(clientID string, name JobName, params []string) Job {

	return Job{ID: fmt.Sprintf("%s@%s", clientID, uuid.New().String()), Name: name, Params: params, Timestamp: time.Now().Unix()}
}

func (j Job) String() string {

	return fmt.Sprintf("Job{ID: %s, Name: %s, Params: %v, Timestamp: %d}", j.ID, j.Name, j.Params, j.Timestamp)
}
