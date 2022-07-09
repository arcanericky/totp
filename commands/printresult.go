package commands

import "fmt"

func printResultf(format string, a ...interface{}) (n int, err error) {
	if !collectionFile.useStdio {
		return fmt.Printf(format, a...)
	}

	return 0, nil
}
