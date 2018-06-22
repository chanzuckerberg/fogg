package templates

import "github.com/gobuffalo/packr"

type T struct {
	Account packr.Box
	Repo    packr.Box
}

var Templates = &T{
	Account: packr.NewBox("account"),
	Repo:    packr.NewBox("repo"),
}
