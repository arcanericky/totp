package cmd

import "fmt"

// func printResultln(a ...interface{}) (n int, err error) {
// 	return fmt.Println(a...)
// }

func printResultf(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf(format, a...)
}
