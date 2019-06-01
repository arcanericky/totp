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
const optionSeed = "seed"

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
		seed, err := cmd.Flags().GetString(optionSeed)
		if err != nil {
			fmt.Println("Error getting seed", err)
			return
		}

		if len(seed) != 0 {
			generateCodeWithSeed(seed)
			return
		}

		if len(args) != 1 {
			fmt.Printf("Need the name of a seed to generate a code.\n\n")
			cmd.Help()
			return
		}

		generateCode(args[0])
	},
}

func generateCodeWithSeed(seed string) {
	code, err := totp.GenerateCode(seed, time.Now())

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error generating code", err)
	} else {
		fmt.Println(code)
	}
}

func generateCode(name string) {
	s, err := api.NewCollectionWithFile(defaultCollectionFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading settings", err)
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
	rootCmd.Flags().StringP(optionSeed, "s", "", "TOTP seed value")
	rootCmd.PersistentFlags().StringP(optionFile, "f", "", "seed collection file")
	rootCmd.SetUsageTemplate(strings.Replace(rootCmd.UsageTemplate(), "{{.UseLine}}", "{{.UseLine}}\n  {{.CommandPath}} [seed name]", 1))
}
