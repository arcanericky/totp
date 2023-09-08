package commands

import (
	"fmt"
	"io"
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
	optionQr       = "qrcode"
	optionSecret   = "secret"
	optionStdio    = "stdio"
	optionTime     = "time"
	optionYes      = "yes"
)

type generateCodesAPI func(time.Duration, time.Duration, time.Duration, func(time.Duration), string, string)

type runVars struct {
	secret     string
	backward   time.Duration
	forward    time.Duration
	timeString string
	follow     bool
	useStdio   bool
	cfgFile    string
	qr         bool
}

var (
	generateCodesService generateCodesAPI
	exitVal              int = 0
)

func getSecretNamesForCompletion(toComplete string) []string {
	var (
		secretNames []string
		err         error
		c           *api.Collection
	)

	c, err = collectionFile.loader()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading collection:", err)
	} else {
		secrets := c.GetSecrets()
		for _, s := range secrets {
			if strings.HasPrefix(s.Name, toComplete) {
				secretNames = append(secretNames, s.Name)
			}
		}
	}

	return secretNames
}

func generateCode(writer io.Writer, name string, secret string, t time.Time) error {
	const errGen = "Error generating code:"
	if len(secret) != 0 {
		code, err := totp.GenerateCode(secret, t)
		if err != nil {
			fmt.Fprintln(os.Stderr, errGen, err)
			return err
		}

		fmt.Fprintln(writer, code)

		return nil
	}

	c, err := collectionFile.loader()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading collection:", err)
		return err
	}

	code, err := c.GenerateCodeWithTime(name, t)
	if err != nil {
		fmt.Fprintln(os.Stderr, errGen, err)
		return err
	}

	fmt.Fprintln(writer, code)

	return nil
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
			if err := generateCode(os.Stdout, secretName, secret, time.Now().Add(timeOffset)); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return true
			}
			return false
		})
}

func run(cmd *cobra.Command, args []string, cfg runVars) int {
	secretLen := len(cfg.secret)
	argsLen := len(args)

	errMsg := ""
	switch {
	// No secret, no secret name
	case secretLen == 0 && argsLen == 0:
		errMsg = "Secret name or secret is required."
	// Secret given but additional arguments were also given
	case !cfg.qr && secretLen > 0 && argsLen > 0:
		errMsg = "Secret was given so additional arguments are not needed."
	// No secret given and too many args
	case secretLen == 0 && argsLen > 1:
		errMsg = "Too many arguments. Only one secret name is required."
	}

	if errMsg != "" {
		fmt.Fprintf(os.Stderr, errMsg+"\n\n")
		if err := cmd.Help(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		return 1
	}

	if cfg.qr {
		secretName := ""
		if argsLen == 1 {
			secretName = args[0]
		}

		if err := qrCode(os.Stdout, secretName, cfg.secret); err != nil {
			return 1
		}

		return 0
	}

	// Override if time was given
	var (
		codeTime time.Time
		err      error
	)
	if len(cfg.timeString) > 0 {
		codeTime, err = time.Parse(time.RFC3339, cfg.timeString)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing the time option:", err)
			return 1
		}
	} else {
		codeTime = time.Now()
		// codeOffset is 0
	}

	// Load the secret name
	secretName := ""
	if argsLen == 1 {
		secretName = args[0]
	}

	// If here then a stored shared secret is wanted
	if err := generateCode(os.Stdout, secretName, cfg.secret, codeTime.Add(cfg.forward-cfg.backward)); err != nil {
		// generateCode will output error text
		return 1
	}

	if cfg.follow {
		generateCodesService(time.Until(codeTime)-cfg.backward+cfg.forward, 0, 30*time.Second, time.Sleep, secretName, cfg.secret)
	}

	return 0
}

func validArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return getSecretNamesForCompletion(toComplete), cobra.ShellCompDirectiveNoFileComp
}

func getRootCmd() *cobra.Command {
	var cfg runVars

	var cobraCmd = &cobra.Command{
		Use:   "totp",
		Short: "TOTP Generator",
		Long:  `TOTP Generator`,
		Args:  cobra.ArbitraryArgs,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Flags().Changed(optionFile) {
				collectionFile.filename = cfg.cfgFile
			}

			if cmd.Flags().Lookup(optionStdio) != nil {
				if cfg.useStdio {
					collectionFile.loader = loadCollectionFromStdin
					collectionFile.useStdio = true
				}
			}
		},
		ValidArgsFunction: validArgs,
		Run: func(cmd *cobra.Command, args []string) {
			exitVal = run(cmd, args, cfg)
		},
	}

	var duration time.Duration

	generateCodesService = generateCodes

	cobraCmd.PersistentFlags().StringVarP(&cfg.cfgFile, optionFile, "f", "", "secret collection file")

	cobraCmd.Flags().StringVarP(&cfg.secret, optionSecret, "s", "", "TOTP secret value")
	cobraCmd.Flags().BoolVarP(&cfg.useStdio, optionStdio, "", false, "load with stdin")
	cobraCmd.Flags().StringVarP(&cfg.timeString, optionTime, "", "", "RFC3339 time for TOTP (2019-06-23T20:00:00-05:00)")
	cobraCmd.Flags().DurationVarP(&cfg.backward, optionBackward, "", duration, "move time backward (ex. \"30s\")")
	cobraCmd.Flags().DurationVarP(&cfg.forward, optionForward, "", duration, "move time forward (ex. \"1m\")")
	cobraCmd.Flags().BoolVarP(&cfg.follow, optionFollow, "", false, "continuous output")
	cobraCmd.Flags().BoolVarP(&cfg.qr, optionQr, "", false, "output QR code")

	cobraCmd.SetUsageTemplate(strings.Replace(cobraCmd.UsageTemplate(), "{{.UseLine}}", "{{.UseLine}}\n  {{.CommandPath}} [secret name]", 1))

	cobraCmd.AddCommand(getVersionCmd())
	cobraCmd.AddCommand(getConfigCmd(cobraCmd))

	return cobraCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() int {
	defaults()
	rootCmd := getRootCmd()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		exitVal = 1
	}

	return exitVal
}
