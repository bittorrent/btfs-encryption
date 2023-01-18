package main

import (
	"fmt"
	"os"
)

func main() {
	if err := cmds.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}
