package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func updateSecret(name, value string) int {
	if isReservedCommand(name) {
		fmt.Fprintln(os.Stderr, "The name \""+name+"\" is reserved for the "+name+" command")
		return 1
	}

	// ignore error because file may not exist
	s, _ := collectionFile.loader()

	secret, err := s.UpdateSecret(name, value)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error updating secret:", err)
		return 1
	}

	if err := s.Save(); err != nil {
		fmt.Fprintln(os.Stderr, "Error saving settings:", err)
		return 1
	}

	action := "Updated"
	if secret.DateAdded == secret.DateModified {
		action = "Added"
	}

	if _, err := printResultf("%s secret %s\n", action, name); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}

func getConfigUpdateCmd(rootCmd *cobra.Command) *cobra.Command {
	var cobraCmd = &cobra.Command{
		Use:               "update",
		Aliases:           []string{"add"},
		Short:             "Add or update a secret",
		Long:              `Add or update a secret`,
		ValidArgsFunction: validArgs,
		Run: func(_ *cobra.Command, args []string) {
			if len(args) != 2 {
				fmt.Fprintln(os.Stderr, "Must provide name and secret")
				return
			}

			exitVal = updateSecret(args[0], args[1])
		},
	}

	cobraCmd.Flags().BoolP(optionStdio, "", false, "load with stdin, save with stdout")
	cobraCmd.SetUsageTemplate(strings.Replace(rootCmd.UsageTemplate(), "{{.UseLine}}", "{{.UseLine}}\n  {{.CommandPath}} [secret name] [secret value]", 1))
	return cobraCmd
}
