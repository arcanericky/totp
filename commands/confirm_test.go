package commands

import (
	"bufio"
	"strings"
	"testing"
)

func Test_userConfirm(t *testing.T) {
	type args struct {
		reader *bufio.Reader
		prompt string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "yes",
			args: args{
				reader: bufio.NewReader(strings.NewReader("yes\n")),
				prompt: "",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "y",
			args: args{
				reader: bufio.NewReader(strings.NewReader("y\n")),
				prompt: "",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "no",
			args: args{
				reader: bufio.NewReader(strings.NewReader("no\n")),
				prompt: "",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "n",
			args: args{
				reader: bufio.NewReader(strings.NewReader("n\n")),
				prompt: "",
			},
			want:    false,
			wantErr: false,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userConfirm(tt.args.reader, tt.args.prompt)
			if (err != nil) != tt.wantErr {
				t.Errorf("userConfirm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("userConfirm() = %v, want %v", got, tt.want)
			}
		})
	}
}
