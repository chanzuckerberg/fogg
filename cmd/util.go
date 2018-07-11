package cmd

import (
	"fmt"
	"os"

	git "gopkg.in/src-d/go-git.v4"
)

func openGitOrExit(pwd string) *git.Repository {
	g, err := git.PlainOpen(pwd)
	if err != nil {
		// assuming this means no repository
		fmt.Println("fogg must be run from the root of a git repo")
		os.Exit(1)
	}
	return g
}
