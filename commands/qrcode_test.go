package commands

import (
	"bytes"
	"encoding/base64"
	"math/rand"
	"os"
	"testing"
)

var testQrCode = `4paI4paI4paI4paI4paI4paI4paI4paI4paI4paI4paI4paI4paI4paI4paI4paI4paI4paI4paI
4paI4paI4paI4paI4paI4paI4paI4paI4paI4paI4paI4paI4paI4paI4paI4paI4paI4paI4paI
4paI4paI4paICuKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKW
iOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKW
iOKWiOKWiOKWiOKWiOKWiOKWiArilojilojilojilogg4paE4paE4paE4paE4paEIOKWiOKWhOKW
hOKWgOKWiOKWiCDiloDiloQg4paI4paE4paA4paI4paA4paE4paE4paA4paIIOKWhOKWhOKWhOKW
hOKWhCDilojilojilojilogK4paI4paI4paI4paIIOKWiCAgIOKWiCDilojiloDiloAg4paI4paA
4paI4paEIOKWgOKWhOKWgOKWiOKWgOKWhOKWiCAg4paIIOKWiCAgIOKWiCDilojilojilojilogK
4paI4paI4paI4paIIOKWiOKWhOKWhOKWhOKWiCDilogg4paEIOKWhOKWgCDilojiloDiloTiloTi
loTiloDilojilojiloTiloTilojilogg4paI4paE4paE4paE4paIIOKWiOKWiOKWiOKWiAriloji
lojilojilojiloTiloTiloTiloTiloTiloTiloTilogg4paAIOKWiCDilojiloTilojiloTiloAg
4paA4paE4paAIOKWiOKWhOKWiOKWhOKWhOKWhOKWhOKWhOKWhOKWhOKWiOKWiOKWiOKWiAriloji
lojilojilojiloTilojiloDilogg4paI4paE4paEIOKWiOKWgOKWiOKWiOKWiOKWhCDiloAg4paA
4paE4paE4paE4paI4paA4paI4paE4paE4paE4paEIOKWgOKWiCDilojilojilojilogK4paI4paI
4paI4paIIOKWgCAgIOKWgOKWhOKWgOKWiOKWiOKWiOKWgOKWhCDiloTiloQg4paA4paE4paA4paA
IOKWiOKWhCDilojiloAg4paE4paE4paIIOKWiOKWiOKWiOKWiOKWiArilojilojilojilojiloDi
loTilojiloTilogg4paE4paA4paE4paA4paIIOKWhOKWhOKWgOKWiOKWhOKWhOKWgOKWiOKWiCDi
loQg4paI4paAICDilojiloTiloDilojiloTilojilojilojilogK4paI4paI4paI4paI4paIIOKW
hCDiloDiloTiloTilojiloQg4paE4paE4paI4paIICDiloAgIOKWiCDiloDilogg4paE4paIIOKW
gOKWgOKWiOKWhCDiloDilojilojilojilogK4paI4paI4paI4paI4paI4paE4paA4paI4paE4paA
4paE4paI4paA4paA4paE4paIIOKWiOKWgOKWgCDiloTiloTiloTiloTiloQgIOKWhCDiloTiloDi
loTiloQg4paA4paA4paI4paI4paI4paICuKWiOKWiOKWiOKWiOKWiOKWgOKWiCAg4paE4paE4paA
4paEIOKWgOKWhCDilojiloDiloTiloAg4paA4paI4paAIOKWiCDiloAg4paI4paEIOKWhOKWiOKW
hOKWiOKWiOKWiOKWiOKWiArilojilojilojilogg4paE4paA4paE4paI4paA4paE4paI4paA4paA
4paIIOKWiOKWhOKWiOKWiOKWgCDiloDiloTiloTiloQg4paA4paE4paIIOKWgOKWiOKWgOKWgOKW
gOKWhOKWiOKWiOKWiOKWiArilojilojilojilojilojilojilojiloTiloTiloTiloTiloTilogg
4paE4paE4paE4paA4paA4paI4paE4paAIOKWiOKWgCDiloDiloTiloAgIOKWiCDiloTilogg4paI
4paI4paI4paI4paICuKWiOKWiOKWiOKWiOKWhOKWhOKWiOKWiOKWiOKWiOKWhOKWiCAgIOKWiOKW
gCDilogg4paAIOKWiOKWhOKWiOKWhOKWhOKWgCDiloTiloTiloQg4paEIOKWhOKWgOKWiOKWiOKW
iOKWiArilojilojilojilogg4paE4paE4paE4paE4paEIOKWiOKWhOKWiOKWiCAg4paI4paAIOKW
iCAg4paA4paE4paAIOKWiCDilojiloTilogg4paI4paIIOKWiOKWiOKWiOKWiOKWiArilojiloji
lojilogg4paIICAg4paIIOKWiOKWhOKWiCDiloTilogg4paI4paI4paI4paE4paA4paE4paI4paE
4paA4paA4paEICAgIOKWhOKWgOKWiOKWiOKWiOKWiOKWiOKWiArilojilojilojilogg4paI4paE
4paE4paE4paIIOKWiOKWiOKWhOKWiOKWiCDiloAgIOKWiOKWgOKWgOKWiCAg4paI4paAIOKWgOKW
hOKWiCDilojiloTilojilojilojilojilojilogK4paI4paI4paI4paI4paE4paE4paE4paE4paE
4paE4paE4paI4paE4paI4paE4paE4paE4paE4paE4paI4paI4paE4paI4paI4paI4paE4paE4paI
4paE4paE4paE4paE4paI4paI4paE4paI4paE4paI4paI4paI4paICuKWiOKWiOKWiOKWiOKWiOKW
iOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKW
iOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiOKWiAriloDiloDi
loDiloDiloDiloDiloDiloDiloDiloDiloDiloDiloDiloDiloDiloDiloDiloDiloDiloDiloDi
loDiloDiloDiloDiloDiloDiloDiloDiloDiloDiloDiloDiloDiloDiloDiloDiloDiloDiloDi
loAK`

func RandStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func Test_getQrString(t *testing.T) {
	type args struct {
		name   string
		secret string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success",
			args: args{
				name:   "testname",
				secret: "testsecret",
			},
			want: "otpauth://totp/testname?secret=TESTSECRET&issuer=testname",
		},
		{
			name: "success",
			args: args{
				name:   "testname",
				secret: "TESTSECRET",
			},
			want: "otpauth://totp/testname?secret=TESTSECRET&issuer=testname",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getQrString(tt.args.name, tt.args.secret); got != tt.want {
				t.Errorf("getQrString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_qrCode(t *testing.T) {
	collectionFile.filename = "testcollection.json"
	createTestData(t)
	collectionFile.loader = loadCollectionFromDefaultFile
	decoded, _ := base64.StdEncoding.DecodeString(testQrCode)
	type args struct {
		name   string
		secret string
	}
	tests := []struct {
		name       string
		filename   string
		args       args
		wantWriter string
		wantErr    bool
	}{
		{
			name:     "name and secret",
			filename: "testcollection.json",
			args: args{
				name:   "testname",
				secret: "testsecret",
			},
			wantWriter: string(decoded),
			wantErr:    false,
		},
		{
			name:     "name from collection",
			filename: "testcollection.json",
			args: args{
				name:   "testname",
				secret: "",
			},
			wantWriter: string(decoded),
			wantErr:    false,
		},
		{
			name:     "empty name",
			filename: "testcollection.json",
			args: args{
				name:   "",
				secret: "testsecret",
			},
			wantWriter: "",
			wantErr:    true,
		},
		{
			name:     "no collecton entry",
			filename: "testcollection.json",
			args: args{
				name:   "invalidname",
				secret: "",
			},
			wantWriter: "",
			wantErr:    true,
		},
		{
			name:     "no collecton file",
			filename: "invalidcollectionfile.json",
			args: args{
				name:   "invalidname",
				secret: "",
			},
			wantWriter: "",
			wantErr:    true,
		},
		{
			name:     "invalid secret",
			filename: "testcollection.json",
			args: args{
				name:   "testname",
				secret: "seed0",
			},
			wantWriter: "",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collectionFile.filename = tt.filename
			writer := &bytes.Buffer{}
			if err := qrCode(writer, tt.args.name, tt.args.secret); (err != nil) != tt.wantErr {
				t.Errorf("qrCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("qrCode() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
	os.Remove(collectionFile.filename)
}

func Test_outputQrCode(t *testing.T) {
	decoded, _ := base64.StdEncoding.DecodeString(testQrCode)

	type args struct {
		name   string
		secret string
	}
	tests := []struct {
		name       string
		args       args
		wantWriter string
		wantErr    bool
	}{
		{
			name: "valid name and secret",
			args: args{
				name:   "testname",
				secret: "TESTSECRET",
			},
			wantWriter: string(decoded),
			wantErr:    false,
		},
		{
			name: "valid name and lowercase secret",
			args: args{
				name:   "testname",
				secret: "testsecret",
			},
			wantWriter: string(decoded),
			wantErr:    false,
		},
		{
			name: "data too long",
			args: args{
				name:   RandStringBytes(2048),
				secret: RandStringBytes(2048),
			},
			wantWriter: "",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			if err := outputQrCode(writer, tt.args.name, tt.args.secret); (err != nil) != tt.wantErr {
				t.Errorf("outputQrCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("outputQrCode() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}
