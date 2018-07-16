package util

import log "github.com/sirupsen/logrus"

func Dump(foo interface{}) {
	log.Printf("%#v\n", foo)
}
