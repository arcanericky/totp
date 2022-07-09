package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func renameSecret(source, target string) {
	if isReservedCommand(target) {
		fmt.Fprintln(os.Stderr, "The name \""+target+"\" is reserved for the "+target+" command")
		return
	}

	s, _ := collectionFile.loader()
	if _, err := s.RenameSecret(source, target); err != nil {
		fmt.Fprintln(os.Stderr, "Error renaming secret:", err)
		return
	}

	if err := s.Save(); err != nil {
		fmt.Fprintln(os.Stderr, "Error saving settings:", err)
		return
	}

	if _, err := printResultf("Renamed secret %s to %s\n", source, target); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}

func getConfigRenameCmd(rootCmd *cobra.Command) *cobra.Command {
	var cobraCmd = &cobra.Command{
		Use:     "rename",
		Aliases: []string{"ren", "mv"},
		Short:   "Rename a secret",
		Long:    `Rename a secret`,
		Run: func(_ *cobra.Command, args []string) {
			if len(args) != 2 {
				fmt.Fprintln(os.Stderr, "Must provide source and target.")
				return
			}

			renameSecret(args[0], args[1])
		},
	}

	cobraCmd.Flags().BoolP(optionStdio, "", false, "load with stdin, save with stdout")
	cobraCmd.SetUsageTemplate(strings.Replace(rootCmd.UsageTemplate(), "{{.UseLine}}", "{{.UseLine}}\n  {{.CommandPath}} [old secret name] [new secret name]", 1))
	return cobraCmd
}
