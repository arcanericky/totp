package commands

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/arcanericky/totp"
	"github.com/spf13/cobra"
)

func titleLine(len int) string {
	var builder strings.Builder
	builder.Grow(len)

	for i := 0; i < len; i++ {
		builder.WriteString("-")
	}

	return builder.String()
}

func listSecretNames(writer io.Writer, secrets []totp.Secret) {
	for _, s := range secrets {
		fmt.Fprintln(writer, s.Name)
	}
}

func listInfo(writer io.Writer, secrets []totp.Secret, all bool) {
	const (
		nameTitle         = "Name"
		secretTitle       = "Secret"
		addedDateTitle    = "Date Added"
		modifiedDateTitle = "Date Modified"
	)

	maxNameLen := len(nameTitle)
	maxSecretLen := len(secretTitle)
	for _, s := range secrets {
		nameLen := len(s.Name)
		if nameLen > maxNameLen {
			maxNameLen = nameLen
		}

		secretLen := len(s.Value)
		if secretLen > maxSecretLen {
			maxSecretLen = secretLen
		}
	}

	const timeFormat = "Jan _2 2006 15:04:05"
	timeFormatLen := len(timeFormat)
	timeFormatLine := titleLine(len(timeFormat))

	if all {
		fmt.Fprintf(writer, "%-*s %-*s %-*s %-*s\n",
			maxNameLen, nameTitle,
			maxSecretLen, secretTitle,
			timeFormatLen, addedDateTitle,
			timeFormatLen, modifiedDateTitle)
		fmt.Fprintf(writer, "%s %s %s %s\n", titleLine(maxNameLen), titleLine(maxSecretLen), timeFormatLine, timeFormatLine)
		for _, s := range secrets {
			fmt.Fprintf(writer, "%-*s %-*s %s %s\n", maxNameLen, s.Name, maxSecretLen, s.Value, s.DateAdded.Format(timeFormat), s.DateModified.Format(timeFormat))
		}
		return
	}

	fmt.Fprintf(writer, "%-*s %-*s %-*s\n",
		maxNameLen, nameTitle,
		timeFormatLen, addedDateTitle,
		timeFormatLen, modifiedDateTitle)
	fmt.Fprintf(writer, "%s %s %s\n", titleLine(maxNameLen), timeFormatLine, timeFormatLine)
	for _, s := range secrets {
		fmt.Fprintf(writer, "%-*s %s %s\n", maxNameLen, s.Name, s.DateAdded.Format(timeFormat), s.DateModified.Format(timeFormat))
	}
}

func listSecrets(names, all bool) int {
	c, err := collectionFile.loader()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading collection", err)
		return 1
	}

	secrets := c.GetSecrets()
	sort.Slice(secrets, func(i, j int) bool {
		return secrets[i].Name < secrets[j].Name
	})

	if names {
		listSecretNames(os.Stdout, secrets)
	} else {
		listInfo(os.Stdout, secrets, all)
	}

	return 0
}

func getConfigListCmd() *cobra.Command {
	var (
		names    bool
		all      bool
		cobraCmd = &cobra.Command{
			Use:     "list",
			Aliases: []string{"ls", "l"},
			Short:   "List secrets",
			Long:    `List secrets`,
			Run: func(listCmd *cobra.Command, _ []string) {
				if names && all {
					fmt.Fprintln(os.Stderr, "Only one of --names or --all can be used.")
					exitVal = 1
					return
				}

				exitVal = listSecrets(names, all)
			},
		}
	)

	cobraCmd.Flags().BoolVarP(&names, "names", "n", false, "list only secret names")
	cobraCmd.Flags().BoolVarP(&all, "all", "a", false, "list all secret info")

	cobraCmd.Flags().BoolP(optionStdio, "", false, "load data from stdin")

	return cobraCmd
}
