package templates

import "github.com/gobuffalo/packr"

type T struct {
	Account   packr.Box
	Component packr.Box
	Env       packr.Box
	Repo      packr.Box
}

var Templates = &T{
	Account:   packr.NewBox("account"),
	Component: packr.NewBox("component"),
	Env:       packr.NewBox("env"),
	Repo:      packr.NewBox("repo"),
}
