package cmd

import (
	"fmt"
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

func listSecretNames(secrets []totp.Secret) {
	for _, s := range secrets {
		fmt.Println(s.Name)
	}
}

func listAllInfo(secrets []totp.Secret) {
	nameTitle := "Name"
	secretTitle := "Secret"
	addedDateTitle := "Date Added"
	modifiedDateTitle := "Date Modified"

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

	timeFormat := "Jan _2 2006 15:04:05"
	timeFormatLen := len(timeFormat)
	timeFormatLine := titleLine(len(timeFormat))
	fmt.Printf("%-*s %-*s %-*s %-*s\n",
		maxNameLen, nameTitle,
		maxSecretLen, secretTitle,
		timeFormatLen, addedDateTitle,
		timeFormatLen, modifiedDateTitle)
	fmt.Printf("%s %s %s %s\n", titleLine(maxNameLen), titleLine(maxSecretLen), timeFormatLine, timeFormatLine)
	for _, s := range secrets {
		fmt.Printf("%-*s %-*s %s %s\n", maxNameLen, s.Name, maxSecretLen, s.Value, s.DateAdded.Format(timeFormat), s.DateModified.Format(timeFormat))
	}
}

func listSecrets(names bool) {
	c, err := collectionFile.loader()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading collection", err)
		return
	}

	secrets := c.GetSecrets()
	sort.Slice(secrets, func(i, j int) bool {
		return secrets[i].Name < secrets[j].Name
	})

	if names {
		listSecretNames(secrets)
	} else {
		listAllInfo(secrets)
	}
}

func getConfigListCmd() *cobra.Command {
	var names bool

	var cobraCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "l"},
		Short:   "List secrets",
		Long:    `List secrets`,
		Run: func(listCmd *cobra.Command, _ []string) {
			listSecrets(names)
		},
	}

	cobraCmd.Flags().BoolVarP(&names, "names", "n", false, "list only secret names")
	cobraCmd.Flags().BoolP(optionStdio, "", false, "load data from stdin")

	return cobraCmd
}
