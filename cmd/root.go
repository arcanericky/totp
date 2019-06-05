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
const optionStdio = "stdio"
const optionTime = "time"

var rootCmd = &cobra.Command{
	Use:   "totp",
	Short: "TOTP Generator",
	Long:  `TOTP Generator`,
	Args:  cobra.ArbitraryArgs,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed(optionFile) {
			cfgFile, err := cmd.Flags().GetString(optionFile)
			if err != nil {
				fmt.Println("Error processing collection file option", err)
				return
			}

			collectionFile.filename = cfgFile
		}

		if cmd.Flags().Lookup(optionStdio) != nil {
			useStdio, err := cmd.Flags().GetBool(optionStdio)
			if err != nil {
				fmt.Println("Error processing stdio option", err)
				return
			}

			if useStdio == true {
				collectionFile.loader = loadCollectionFromStdin
				collectionFile.useStdio = true
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Process the secret option
		secret, err := cmd.Flags().GetString(optionSecret)
		if err != nil {
			fmt.Println("Error getting secret", err)
			return
		}

		// Process the time option
		timeString, err := cmd.Flags().GetString(optionTime)
		if err != nil {
			fmt.Println("Error processing time option", err)
			return
		}

		codeTime := time.Now()

		if len(timeString) > 0 {
			codeTime, err = time.Parse(time.RFC3339, timeString)
			if err != nil {
				fmt.Println("Error parsing the time option", err)
				return
			}
		}

		// Providing a secret name overrides the --secret option but
		// should probably generate an error if both are given
		if len(secret) != 0 {
			generateCodeWithSecret(secret, codeTime)
			return
		}

		// If no secret or no secret name given, show help
		if len(args) != 1 {
			fmt.Fprintf(os.Stderr, "Need the name of a secret to generate a code.\n\n")
			cmd.Help()
			return
		}

		// If here then a stored shared secret is wanted
		generateCode(args[0], codeTime)
	},
}

func generateCodeWithSecret(secret string, t time.Time) {
	code, err := totp.GenerateCode(secret, t)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error generating code", err)
	} else {
		fmt.Println(code)
	}
}

func generateCode(name string, t time.Time) {
	var s *api.Collection
	var err error

	s, err = collectionFile.loader()

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading collection", err)
	} else {
		code, err := s.GenerateCodeWithTime(name, t)

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
	rootCmd.PersistentFlags().StringP(optionFile, "f", "", "secret collection file")

	rootCmd.Flags().StringP(optionSecret, "s", "", "TOTP secret value")
	rootCmd.Flags().BoolP(optionStdio, "", false, "load with stdin")
	rootCmd.Flags().StringP(optionTime, "", "", "RFC3339 time for TOTP (2006-01-02T15:04:05Z07:00)")
	rootCmd.SetUsageTemplate(strings.Replace(rootCmd.UsageTemplate(), "{{.UseLine}}", "{{.UseLine}}\n  {{.CommandPath}} [secret name]", 1))
}
