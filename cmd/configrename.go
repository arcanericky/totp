package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var configRenameCmd = &cobra.Command{
	Use:     "rename",
	Aliases: []string{"ren", "mv"},
	Short:   "Rename a secret",
	Long:    `Rename a secret`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Fprintln(os.Stderr, "Must provide source and target.")
			return
		}

		renameSecret(args[0], args[1])
	},
}

func renameSecret(source, target string) {
	if isReservedCommand(target) {
		fmt.Fprintln(os.Stderr, "The name \""+target+"\" is reserved for the "+target+" command")
		return
	}

	s, _ := collectionFile.loader()
	_, err := s.RenameSecret(source, target)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error renaming secret:", err)
	} else {
		s.Save()
		printResultf("Renamed secret %s to %s\n", source, target)
	}
}

func init() {
	configCmd.AddCommand(configRenameCmd)
	configRenameCmd.Flags().BoolP(optionStdio, "", false, "load with stdin, save with stdout")
	configRenameCmd.SetUsageTemplate(strings.Replace(rootCmd.UsageTemplate(), "{{.UseLine}}", "{{.UseLine}}\n  {{.CommandPath}} [old secret name] [new secret name]", 1))
}
