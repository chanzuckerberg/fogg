package util

import "github.com/sirupsen/logrus"

func Dump(foo interface{}) {
	logrus.Debugf("%#v\n", foo)
}
