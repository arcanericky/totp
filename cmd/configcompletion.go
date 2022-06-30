package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func getConfigCompletionCmd(rootCmd *cobra.Command) *cobra.Command {
	var cobraCmd = &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		Long: `To load completion run

. <(totp config completion)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(totp config completion)
`,
		Run: func(_ *cobra.Command, _ []string) {
			if err := rootCmd.GenBashCompletion(os.Stdout); err != nil {
				fmt.Fprintf(os.Stderr, "Error generating completion script: %s\n", err)
			}
		},
	}

	return cobraCmd
}
