package templates

import "github.com/gobuffalo/packr"

type T struct {
	Repo packr.Box
}

var Templates = &T{
	Repo: packr.NewBox("repo"),
}
