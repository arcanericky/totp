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

const (
	optionBackward = "backward"
	optionFile     = "file"
	optionFollow   = "follow"
	optionForward  = "forward"
	optionSecret   = "secret"
	optionStdio    = "stdio"
	optionTime     = "time"
)

type generateCodesAPI func(time.Duration, time.Duration, time.Duration, func(time.Duration), string, string)

var generateCodesService generateCodesAPI

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

		// Process the backward option
		backward, err := cmd.Flags().GetDuration(optionBackward)
		if err != nil {
			fmt.Println("Error processing backward option", err)
			return
		}

		// Process the forward option
		forward, err := cmd.Flags().GetDuration(optionForward)
		if err != nil {
			fmt.Println("Error processing forward option", err)
			return
		}

		// Process the time option
		timeString, err := cmd.Flags().GetString(optionTime)
		if err != nil {
			fmt.Println("Error processing time option", err)
			return
		}

		follow, err := cmd.Flags().GetBool(optionFollow)
		if err != nil {
			fmt.Println("Error processing follow option", err)
			return
		}

		var codeTime time.Time

		// Override if time was given
		if len(timeString) > 0 {
			codeTime, err = time.Parse(time.RFC3339, timeString)
			if err != nil {
				fmt.Println("Error parsing the time option", err)
				return
			}
		} else {
			codeTime = time.Now()
			// codeOffset is 0
		}

		secretLen := len(secret)
		argsLen := len(args)

		// No secret, no secret name
		if secretLen == 0 && argsLen == 0 {
			fmt.Fprintf(os.Stderr, "Secret name or secret is required.\n\n")
			cmd.Help()

			return
		}

		// Secret given but additional arguments were also given
		if secretLen > 0 && argsLen > 0 {
			fmt.Fprintf(os.Stderr, "Secret was given so additional arguments are not needed.\n\n")
			cmd.Help()

			return
		}

		// No secret given and too many args
		if secretLen == 0 && argsLen > 1 {
			fmt.Fprintf(os.Stderr, "Too many arguments. Only one secret name is required.\n\n")
			cmd.Help()

			return
		}

		// Load the secret name
		secretName := ""
		if argsLen == 1 {
			secretName = args[0]
		}

		// If here then a stored shared secret is wanted
		generateCode(secretName, secret, codeTime.Add(forward-backward))

		if follow {
			generateCodesService(codeTime.Sub(time.Now())-backward+forward, 0, 30*time.Second, time.Sleep, secretName, secret)
		}
	},
}

func generateCode(name string, secret string, t time.Time) error {
	var code string
	var err error
	var c *api.Collection

	if len(secret) != 0 {
		code, err = totp.GenerateCode(secret, t)
	} else {
		c, err = collectionFile.loader()

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error loading collection", err)
		} else {
			code, err = c.GenerateCodeWithTime(name, t)

			if err != nil {
				fmt.Fprintln(os.Stderr, "Error generating code", err)
			}
		}
	}

	if err == nil {
		fmt.Println(code)
	}

	return err
}

func durationToNextInterval(now time.Time) time.Duration {
	var sleepSeconds int

	s := now.Second()
	switch {
	case s == 0, s < 30:
		sleepSeconds = 30 - s
	case s >= 30:
		sleepSeconds = 60 - s
	}

	return time.Duration(sleepSeconds)*time.Second -
		time.Duration(now.Nanosecond())*time.Nanosecond

}

func callOnInterval(runtime time.Duration, interval time.Duration, exec func() bool) {
	stopper := make(chan bool)

	if runtime > 0 {
		go func() {
			time.Sleep(runtime)
			stopper <- true
		}()
	}

	if exec != nil && exec() {
		return
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-stopper:
			return
		case <-ticker.C:
			if exec != nil && exec() {
				return
			}
		}
	}
}

func generateCodes(timeOffset time.Duration, durationToRun time.Duration, intervalTime time.Duration, sleep func(time.Duration), secretName, secret string) {
	sleep(durationToNextInterval(time.Now().Add(timeOffset)) + 10*time.Millisecond)

	callOnInterval(durationToRun, intervalTime,
		func() bool {
			generateCode(secretName, secret, time.Now().Add(timeOffset))
			return false
		})
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
	var duration time.Duration

	generateCodesService = generateCodes

	rootCmd.PersistentFlags().StringP(optionFile, "f", "", "secret collection file")

	rootCmd.Flags().StringP(optionSecret, "s", "", "TOTP secret value")
	rootCmd.Flags().BoolP(optionStdio, "", false, "load with stdin")
	rootCmd.Flags().StringP(optionTime, "", "", "RFC3339 time for TOTP (2019-06-23T20:00:00-05:00)")
	rootCmd.Flags().DurationP(optionBackward, "", duration, "move time backward (ex. \"30s\")")
	rootCmd.Flags().DurationP(optionForward, "", duration, "move time forward (ex. \"1m\")")
	rootCmd.Flags().BoolP(optionFollow, "", false, "continuous output")

	rootCmd.SetUsageTemplate(strings.Replace(rootCmd.UsageTemplate(), "{{.UseLine}}", "{{.UseLine}}\n  {{.CommandPath}} [secret name]", 1))
}
