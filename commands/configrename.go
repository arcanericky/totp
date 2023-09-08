package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func renameSecret(source, target string) int {
	if isReservedCommand(target) {
		fmt.Fprintln(os.Stderr, "The name \""+target+"\" is reserved for the "+target+" command")
		return 1
	}

	s, _ := collectionFile.loader()
	if _, err := s.RenameSecret(source, target); err != nil {
		fmt.Fprintln(os.Stderr, "Error renaming secret:", err)
		return 1
	}

	if err := s.Save(); err != nil {
		fmt.Fprintln(os.Stderr, "Error saving settings:", err)
		return 1
	}

	if _, err := printResultf("Renamed secret %s to %s\n", source, target); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}

func getConfigRenameCmd(rootCmd *cobra.Command) *cobra.Command {
	var cobraCmd = &cobra.Command{
		Use:               "rename",
		Aliases:           []string{"ren", "mv"},
		Short:             "Rename a secret",
		Long:              `Rename a secret`,
		ValidArgsFunction: validArgs,
		Run: func(_ *cobra.Command, args []string) {
			if len(args) != 2 {
				fmt.Fprintln(os.Stderr, "Must provide source and target.")
				exitVal = 1
				return
			}

			exitVal = renameSecret(args[0], args[1])
		},
	}

	cobraCmd.Flags().BoolP(optionStdio, "", false, "load with stdin, save with stdout")
	cobraCmd.SetUsageTemplate(strings.Replace(rootCmd.UsageTemplate(), "{{.UseLine}}", "{{.UseLine}}\n  {{.CommandPath}} [old secret name] [new secret name]", 1))
	return cobraCmd
}
