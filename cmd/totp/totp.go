package main

import (
	"os"

	"github.com/arcanericky/totp/commands"
)

func main() {
	os.Exit(commands.Execute())
}
