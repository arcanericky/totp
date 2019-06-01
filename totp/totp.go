package main

import (
	"os"

	"github.com/arcanericky/totp/cmd"
)

func main() {
	os.Exit(cmd.Execute())
}
