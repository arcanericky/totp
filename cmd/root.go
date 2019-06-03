package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	api "github.com/arcanericky/totp"
	"github.com/pquerna/otp/totp"
	"github.com/spf13/cobra"
)

const optionFile = "file"
const optionSecret = "secret"
const optionStdin = "stdin"

var rootCmd = &cobra.Command{
	Use:   "totp",
	Short: "TOTP Generator",
	Long:  `TOTP Generator`,
	Args:  cobra.ArbitraryArgs,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed(optionFile) {
			cfgFile, err := cmd.Flags().GetString(optionFile)
			if err != nil {
				fmt.Println("Error process collection file option", err)
				return
			}

			defaultCollectionFile = cfgFile
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		secret, err := cmd.Flags().GetString(optionSecret)
		if err != nil {
			fmt.Println("Error getting secret", err)
			return
		}

		if len(secret) != 0 {
			generateCodeWithSecret(secret)
			return
		}

		if len(args) != 1 {
			fmt.Printf("Need the name of a secret to generate a code.\n\n")
			cmd.Help()
			return
		}

		generateCode(args[0])
	},
}

func generateCodeWithSecret(secret string) {
	code, err := totp.GenerateCode(secret, time.Now())

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error generating code", err)
	} else {
		fmt.Println(code)
	}
}

func generateCode(name string) {
	var s *api.Collection
	var err error

	s, err = api.NewCollectionWithFile(defaultCollectionFile)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading collection", err)
	} else {
		code, err := s.GenerateCode(name)

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error generating code", err)
		} else {
			fmt.Println(code)
		}
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() int {
	retVal := 0

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		retVal = 1
	}

	return retVal
}

func init() {
	rootCmd.Flags().StringP(optionSecret, "s", "", "TOTP secret value")
	rootCmd.PersistentFlags().StringP(optionFile, "f", "", "secret collection file")

	rootCmd.SetUsageTemplate(strings.Replace(rootCmd.UsageTemplate(), "{{.UseLine}}", "{{.UseLine}}\n  {{.CommandPath}} [secret name]", 1))
}
