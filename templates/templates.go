package templates

import "github.com/gobuffalo/packr"

type T struct {
	Account   packr.Box
	Component packr.Box
	Env       packr.Box
	Global    packr.Box
	Repo      packr.Box
}

var Templates = &T{
	Account:   packr.NewBox("account"),
	Component: packr.NewBox("component"),
	Env:       packr.NewBox("env"),
	Global:    packr.NewBox("global"),
	Repo:      packr.NewBox("repo"),
}
