package main

import (
	"os"

	"github.com/kardianos/service"
)

func main() {

	svcConfig := &service.Config{
		Name:        "Gowin",
		DisplayName: "Gowin",
		Description: "Written by @ranjbar-dev",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {

		panic(err)
	}

	// Handle CLI commands
	if len(os.Args) > 1 {

		cmd := os.Args[1]
		err = service.Control(s, cmd)
		if err != nil {

			panic(err)
		}

		return
	}

	err = s.Run()
	if err != nil {

		panic(err)
	}
}
