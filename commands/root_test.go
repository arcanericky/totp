package commands

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func TestRoot(t *testing.T) {
	Execute()
	collectionFile.filename = "testcollection.json"

	secretList := createTestData(t)

	rootCmd := getRootCmd()

	// No parameters
	rootCmd.Run(rootCmd, []string{})

	// Valid entry and secret
	rootCmd.Run(rootCmd, []string{secretList[0].name})

	// Non-existing entry
	rootCmd.Run(rootCmd, []string{"invalidsecret"})

	// Test follow condition
	savedGenerateCodesService := generateCodesService
	generateCodesService = func(time.Duration, time.Duration, time.Duration, func(time.Duration), string, string) {}
	_ = rootCmd.Flags().Lookup(optionFollow).Value.Set("true")
	rootCmd.Run(rootCmd, []string{"name0"})
	generateCodesService = savedGenerateCodesService
	_ = rootCmd.Flags().Lookup(optionFollow).Value.Set("false")

	// Completion
	rootCmd.ValidArgsFunction(rootCmd, []string{}, "na")

	// Completion with args
	rootCmd.ValidArgsFunction(rootCmd, []string{"secret"}, "na")

	// No collections file
	os.Remove(collectionFile.filename)
	rootCmd.Run(rootCmd, []string{secretList[0].name})

	// Completion without collections
	rootCmd.ValidArgsFunction(rootCmd, []string{}, "na")

	// Excessive args
	rootCmd.Run(rootCmd, []string{"secretname", "extraarg"})

	// Provide secret option
	_ = rootCmd.Flags().Set(optionSecret, "seed")
	rootCmd.Run(rootCmd, []string{})

	// Provide invalid secret option
	_ = rootCmd.Flags().Set(optionSecret, "seed1")
	rootCmd.Run(rootCmd, []string{})

	// File option
	_ = rootCmd.Flags().Set(optionFile, collectionFile.filename)
	rootCmd.Flags().Lookup(optionFile).Changed = true
	rootCmd.PersistentPreRun(rootCmd, []string{"secret"})

	// Stdio option
	_ = rootCmd.Flags().Set(optionStdio, "true")
	rootCmd.PersistentPreRun(rootCmd, []string{"secret"})
	collectionFile.loader = loadCollectionFromDefaultFile
	collectionFile.useStdio = false
	_ = rootCmd.Flags().Set(optionStdio, "")

	// Time option
	_ = rootCmd.Flags().Set(optionTime, "2019-06-01T20:00:00-05:00")
	rootCmd.Run(rootCmd, []string{})

	// Give secret and secret name
	_ = rootCmd.Flags().Set(optionSecret, "seed")
	rootCmd.Run(rootCmd, []string{"secretname"})
	_ = rootCmd.Flags().Set(optionSecret, "")

	// Invalid time option
	_ = rootCmd.Flags().Set(optionTime, "invalidtime")
	rootCmd.Run(rootCmd, []string{})
	_ = rootCmd.Flags().Set(optionTime, "")
	os.Remove(collectionFile.filename)
}

func Test_run(t *testing.T) {
	collectionFile.filename = "testcollection.json"
	_ = createTestData(t)
	type args struct {
		cmd  *cobra.Command
		args []string
		cfg  runVars
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "qr code",
			args: args{
				cmd:  &cobra.Command{},
				args: []string{"testname"},
				cfg:  runVars{qr: true},
			},
		},
		{
			name: "secret name or secret required",
			args: args{
				cmd:  &cobra.Command{},
				args: []string{},
				cfg:  runVars{secret: ""},
			},
		},
		{
			name: "secret + additional args",
			args: args{
				cmd:  &cobra.Command{},
				args: []string{"name"},
				cfg: runVars{
					qr:     false,
					secret: "seed",
				},
			},
		},
		{
			name: "too many arguments",
			args: args{
				cmd:  &cobra.Command{},
				args: []string{"name", "extra-arg"},
				cfg: runVars{
					qr:     false,
					secret: "",
				},
			},
		},
		{
			name: "secret from collection",
			args: args{
				cmd:  &cobra.Command{},
				args: []string{"name3"},
				cfg: runVars{
					qr:     false,
					secret: "",
				},
			},
		},
		{
			name: "secret from collection now found",
			args: args{
				cmd:  &cobra.Command{},
				args: []string{"invalidsecretname"},
				cfg: runVars{
					qr:     false,
					secret: "",
				},
			},
		},
		{
			name: "time parse error",
			args: args{
				cmd:  &cobra.Command{},
				args: []string{"testname"},
				cfg: runVars{
					timeString: "invalidtime",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			run(tt.args.cmd, tt.args.args, tt.args.cfg)
		})
	}

	os.Remove(collectionFile.filename)
}

func Test_generateCodes(t *testing.T) {
	var duration time.Duration
	type args struct {
		timeOffset    time.Duration
		durationToRun time.Duration
		intervalTime  time.Duration
		sleep         func(time.Duration)
		secretName    string
		secret        string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "valid seed",
			args: args{
				timeOffset:    duration,
				durationToRun: 2 * time.Millisecond,
				intervalTime:  1 * time.Millisecond,
				sleep:         func(d time.Duration) {},
				secretName:    "",
				secret:        "seed",
			},
		},
		{
			name: "invalid seed",
			args: args{
				timeOffset:    duration,
				durationToRun: 2 * time.Millisecond,
				intervalTime:  1 * time.Millisecond,
				sleep:         func(d time.Duration) {},
				secretName:    "",
				secret:        "invalidseed",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generateCodes(tt.args.timeOffset, tt.args.durationToRun, tt.args.intervalTime, tt.args.sleep, tt.args.secretName, tt.args.secret)
		})
	}
}

func Test_callOnInterval(t *testing.T) {
	execCount := 0
	preAndExecNormal := func() bool { return false }
	preAndExecEarlyExit := func() bool { execCount++; return execCount == 2 }
	type args struct {
		runtime  time.Duration
		interval time.Duration
		exec     func() bool
	}
	tests := []struct {
		name           string
		startExecCount int
		args           args
	}{
		{
			name:           "normal execution",
			startExecCount: 0,
			args: args{
				runtime:  2 * time.Millisecond,
				interval: 1 * time.Millisecond,
				exec:     preAndExecNormal,
			},
		},
		{
			name:           "exit at preExec",
			startExecCount: 1,
			args: args{
				runtime:  2 * time.Millisecond,
				interval: 1 * time.Millisecond,
				exec:     preAndExecNormal,
			},
		},
		{
			name:           "exit at top exec",
			startExecCount: 1,
			args: args{
				runtime:  2 * time.Millisecond,
				interval: 1 * time.Millisecond,
				exec:     preAndExecEarlyExit,
			},
		},
		{
			name:           "exit at loop exec",
			startExecCount: 0,
			args: args{
				runtime:  2 * time.Millisecond,
				interval: 1 * time.Millisecond,
				exec:     preAndExecEarlyExit,
			},
		},
		{
			name:           "no callbacks",
			startExecCount: 0,
			args: args{
				runtime:  2 * time.Millisecond,
				interval: 1 * time.Millisecond,
				exec:     nil,
			},
		},
	}
	for _, tt := range tests {
		execCount = tt.startExecCount
		t.Run(tt.name, func(t *testing.T) {
			callOnInterval(tt.args.runtime, tt.args.interval, tt.args.exec)
		})
	}
}

func Test_durationToNextInterval(t *testing.T) {
	type args struct {
		now string
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "29 seconds",
			args: args{
				now: "2019-06-23T20:00:01-05:00",
			},
			want: time.Duration(29 * time.Second),
		},
		{
			name: "29 seconds 2",
			args: args{
				now: "2019-06-23T20:00:31-05:00",
			},
			want: time.Duration(29 * time.Second),
		},
		{
			name: "28 seconds",
			args: args{
				now: "2019-06-23T20:00:31.001-05:00",
			},
			want: time.Duration(28*time.Second + 999*time.Millisecond),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now, _ := time.Parse(time.RFC3339, tt.args.now)
			if got := durationToNextInterval(now); got != tt.want {
				t.Errorf("durationToNextInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateCode(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2019-06-23T20:00:01-05:00")
	type args struct {
		name   string
		secret string
		t      time.Time
	}
	tests := []struct {
		name       string
		args       args
		wantWriter string
		wantErr    bool
	}{
		{
			name: "valid secret",
			args: args{
				name:   "name",
				secret: "seed",
				t:      now,
			},
			wantWriter: "335072\n",
			wantErr:    false,
		},
		{
			name: "invalid secret",
			args: args{
				name:   "name",
				secret: "seed0",
				t:      now,
			},
			wantWriter: "",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			if err := generateCode(writer, tt.args.name, tt.args.secret, tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("generateCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("generateCode() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func Test_getSecretNamesForCompletion(t *testing.T) {
	collectionFile.filename = "testcollection.json"
	collectionFile.loader = loadCollectionFromDefaultFile
	_ = createTestData(t)
	type args struct {
		toComplete string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "names for completion",
			args: args{
				toComplete: "n",
			},
			want: []string{"name0", "name1", "name2", "name3", "name4"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getSecretNamesForCompletion(tt.args.toComplete)
			for _, want := range tt.want {
				i := 0
				match := false
				for i = range got {
					if want == got[i] {
						match = true
						break
					}
				}
				if !match {
					t.Errorf("want: %s, got: %s", tt.want, got)
				}
			}
		})
	}
	os.Remove(collectionFile.filename)
}
