package main

import (
	"fmt"

	"github.com/kardianos/service"
	"github.com/ranjbar-dev/gowin/srv"
)

// Service configuration
type program struct{}

func (p *program) Start(s service.Service) error {

	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {

	fmt.Println("Stopping service...")
	// Cleanup logic here
	return nil
}

func (p *program) run() {

	srv.StartMain()
}
