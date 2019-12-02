package main

import (
	"github.com/chanzuckerberg/fogg/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	formatter := &logrus.TextFormatter{
		DisableTimestamp: true,
	}
	logrus.SetFormatter(formatter)
	cmd.Execute()
}
