package cmd

import (
	"fmt"
	"os"
)

// fatal prints an error message to stderr and exits with code 1.
func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "错误: "+format+"\n", args...)
	os.Exit(1)
}

