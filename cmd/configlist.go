package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/arcanericky/totp"
	"github.com/spf13/cobra"
)

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List keys",
	Long:  `List keys`,
	Run: func(cmd *cobra.Command, args []string) {
		listKeys(cmd)
	},
}

func titleLine(len int) string {
	var builder strings.Builder
	builder.Grow(len)

	for i := 0; i < len; i++ {
		builder.WriteString("-")
	}

	return builder.String()
}

func listKeyNames(keys []totp.Entry) {
	for _, k := range keys {
		fmt.Println(k.Name)
	}
}

func listAllInfo(keys []totp.Entry) {
	nameTitle := "Name"
	seedTitle := "Seed"
	addedDateTitle := "Date Added"
	modifiedDateTitle := "Date Modified"

	maxNameLen := len(nameTitle)
	maxSeedLen := len(seedTitle)
	for _, k := range keys {
		nameLen := len(k.Name)
		if nameLen > maxNameLen {
			maxNameLen = nameLen
		}

		seedLen := len(k.Seed)
		if seedLen > maxSeedLen {
			maxSeedLen = seedLen
		}
	}

	timeFormat := "Jan _2 2006 15:04:05"
	timeFormatLen := len(timeFormat)
	timeFormatLine := titleLine(len(timeFormat))
	fmt.Printf("%-*s %-*s %-*s %-*s\n",
		maxNameLen, nameTitle,
		maxSeedLen, seedTitle,
		timeFormatLen, addedDateTitle,
		timeFormatLen, modifiedDateTitle)
	fmt.Printf("%s %s %s %s\n", titleLine(maxNameLen), titleLine(maxSeedLen), timeFormatLine, timeFormatLine)
	for _, k := range keys {
		fmt.Printf("%-*s %-*s %s %s\n", maxNameLen, k.Name, maxSeedLen, k.Seed, k.DateAdded.Format(timeFormat), k.DateModified.Format(timeFormat))
	}
}

func listKeys(cmd *cobra.Command) {
	c, err := totp.NewCollectionWithFile(defaultCollectionFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading settings", err)
	} else {
		names, err := cmd.Flags().GetBool("names")
		if err != nil {
			fmt.Println("Error getting names option", err)
			return
		}

		keys := c.GetKeys()

		if names == true {
			listKeyNames(keys)
		} else {
			listAllInfo(keys)
		}
	}
}

func init() {
	configCmd.AddCommand(configListCmd)
	configListCmd.Flags().BoolP("names", "n", false, "list only seed names")
}
