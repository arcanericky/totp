package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/arcanericky/totp"
	"github.com/spf13/cobra"
)

var configRenameCmd = &cobra.Command{
	Use:     "rename",
	Aliases: []string{"ren", "mv"},
	Short:   "Rename a key",
	Long:    `Rename a key`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Fprintln(os.Stderr, "Must provide source and target.")
			return
		}

		renameKey(args[0], args[1])
	},
}

func renameKey(source, target string) {
	if isReservedCommand(target) {
		fmt.Fprintln(os.Stderr, "The name \""+target+"\" is reserved for the "+target+" command")
		return
	}

	s, _ := totp.NewCollectionWithFile(defaultCollectionFile)
	_, err := s.RenameKey(source, target)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error renaming key:", err)
	} else {
		s.Save()
		fmt.Printf("Renamed key %s to %s\n", source, target)
	}
}

func init() {
	configCmd.AddCommand(configRenameCmd)
	configRenameCmd.SetUsageTemplate(strings.Replace(rootCmd.UsageTemplate(), "{{.UseLine}}", "{{.UseLine}}\n  {{.CommandPath}} [old seed name] [new seed name]", 1))
}
