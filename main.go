package main

import (
	"github.com/chanzuckerberg/fogg/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	formatter := &log.TextFormatter{
		DisableTimestamp: true,
	}
	log.SetFormatter(formatter)
	cmd.Execute()
}
