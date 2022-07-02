package commands

import (
	"bufio"
	"fmt"
	"strings"
)

func userConfirm(reader *bufio.Reader, prompt string) (bool, error) {
	for {
		fmt.Printf("%s Continue? [y/n]: ", prompt)

		response, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true, nil
		} else if response == "n" || response == "no" {
			return false, nil
		}
	}
}
