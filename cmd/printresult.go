package cmd

import "fmt"

func printResultf(format string, a ...interface{}) (n int, err error) {
	if collectionFile.useStdio == false {
		return fmt.Printf(format, a...)
	}

	return 0, nil
}
