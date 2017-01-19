package util

import (
	"log"

	git "gopkg.in/src-d/go-git.v4"
)

// Check checks for errors
func Check(err error) {
	if err == nil || err == git.NoErrAlreadyUpToDate {
		return
	}
	log.Fatal(err)
}
