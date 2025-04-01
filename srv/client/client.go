package client

import (
	"github.com/go-resty/resty/v2"
	"github.com/ranjbar-dev/gowin/config"
	"github.com/ranjbar-dev/gowin/types"
)

type Client struct {
	id    string
	jobs  chan types.Job
	host  string
	resty *resty.Client
}

func (c *Client) request() *resty.Request {

	return c.resty.SetDebug(false).SetDisableWarn(true).NewRequest().SetBasicAuth(config.ApiBasicUsername(), config.ApiBasicPassword())
}

func NewClient(id string, host string) *Client {

	return &Client{
		id:    id,
		jobs:  make(chan types.Job, 100),
		host:  host,
		resty: resty.New(),
	}
}

func (c *Client) Start() {

	go c.listenForJobs()

	go c.handleJobs()
}
