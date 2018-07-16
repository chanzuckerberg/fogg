package util

import log "github.com/sirupsen/logrus"

func Dump(foo interface{}) {
	log.Debugf("%#v\n", foo)
}
