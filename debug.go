package prism

import (
	"log"

	"github.com/mk2/prism/env"
)

/*
Debug Utilities
*/

var dbg = debugT(env.Debug)

type debugT bool

func (d debugT) Printf(format string, args ...interface{}) {
	if d {
		log.Printf(format, args...)
	}
}
