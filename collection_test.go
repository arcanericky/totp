package totp

import (
	"bytes"
	"errors"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

type secretItem struct {
	name  string
	value string
}

func createTestData(t *testing.T) []secretItem {
	t.Helper()

	// Create test data
	c, _ := NewCollectionWithFile("testcollection.json")

	// Create some test data
	secretList := []secretItem{
		{name: "name0", value: "seed"},
		{name: "name1", value: "seed"},
		{name: "name2", value: "seedseed"},
		{name: "name3", value: "seed"},
		{name: "name4", value: "seed"},
	}

	for _, i := range secretList {
		_, err := c.UpdateSecret(i.name, i.value)
		if err != nil {
			t.Errorf("Error adding secret %s for test data: %s", i, err)
		}
	}

	_ = c.Save()

	return secretList
}

// func testDeleteSecret(t *testing.T, testSecret string, c *Collection) {
// 	t.Helper()

// 	_, err := c.DeleteSecret(testSecret)
// 	if err != nil {
// 		t.Error("DeleteSecret on valid secret name returned error", err)
// 	}

// 	_, err = c.GetSecret(testSecret)
// 	if err == nil {
// 		t.Error("Error deleting secret", testSecret)
// 	}
// }

// func xTestWriteProtected(t *testing.T) {
// 	secretName := "name"
// 	secretValue := "seed"

// 	updatedValue := "updatedvalue"

// 	s := NewCollection()

// 	secret, err := s.UpdateSecret(secretName, secretValue)
// 	if err != nil {
// 		t.Error("Failed to add secret for test")
// 	}

// 	secret.Name = updatedValue
// 	secret.Value = updatedValue

// 	// Validate secret returned on update can't be changed
// 	secret, _ = s.GetSecret(secretName)
// 	if secret.Name != secretName || secret.Value != secretValue {
// 		t.Error("Internal collection secret can be updated with returned secret from UpdateSecret()")
// 	}
// }

// func xTestSettingsNew(t *testing.T) {
// 	collectionFile := "testcollection.json"

// 	// Test failure on Reader interface
// 	_, err := NewCollectionWithReader(os.Stdout)
// 	if err == nil {
// 		t.Error("New collection should fail with os.Stdout as reader")
// 	}

// 	c := NewCollection()

// 	// Test error on Save with no filename, no writer
// 	err = c.Save()
// 	if err == nil {
// 		t.Error("Save collection with no filename or writer should generate error")
// 	}

// 	// Set filename for remainder of tests
// 	c.SetFilename(collectionFile)

// 	// Create some data
// 	type secretItem struct {
// 		name  string
// 		value string
// 	}

// 	// Create some test data
// 	secretList := []secretItem{
// 		{name: "name0", value: "seed"},
// 		{name: "name1", value: "seed"},
// 		{name: "name2", value: "seed"},
// 		{name: "name3", value: "seed"},
// 		{name: "name4", value: "seed"},
// 	}

// 	for _, i := range secretList {
// 		_, err := c.UpdateSecret(i.name, i.value)
// 		if err != nil {
// 			t.Error("Error updating secret:", err)
// 		}
// 	}

// 	err = c.Save()
// 	if err != nil {
// 		t.Error("Save collection with filename yielded error")
// 	}

// 	// Load test data
// 	c, _ = NewCollectionWithFile(collectionFile)
// 	for _, i := range secretList {
// 		secret, err := c.GetSecret(i.name)
// 		if err == nil {
// 			if secret.Name != i.name || secret.Value != i.value {
// 				t.Error("Loaded secrets don't match saved secrets")
// 			}
// 		} else {
// 			t.Error("Error loading test data:", err)
// 		}
// 	}

// 	// Test GenerateCode() methods
// 	testSecret := secretList[0].name
// 	_, err = c.GenerateCode(testSecret)
// 	if err != nil {
// 		t.Error("Error generating code for secret", testSecret)
// 	}

// 	// Attempt invalid secret retrieval
// 	secret, err := c.GetSecret("invalidsecret")
// 	if err == nil {
// 		t.Error("GetSecret returned success on invalid secret retrieval")
// 	}

// 	newSecret := "deadbeef"

// 	// Update secret with empty name
// 	secret, err = c.UpdateSecret("", newSecret)
// 	if err == nil {
// 		t.Error("UpdateSecret with empty name did not return error")
// 	}

// 	// Update secret with empty value
// 	secret, err = c.UpdateSecret(secretList[0].name, "")
// 	if err == nil {
// 		t.Error("UpdateSecret with empty value did not return error")
// 	}

// 	// Update a secret
// 	testSecret = secretList[0].name
// 	secret, err = c.UpdateSecret(testSecret, newSecret)
// 	if err != nil {
// 		t.Error("Error updating secret", secret, err)
// 	}
// 	if secret.DateAdded == secret.DateModified {
// 		t.Error("Date modified not updated on secret update")
// 	}

// 	secret, err = c.GetSecret(testSecret)
// 	if err != nil || secret.Value != newSecret {
// 		t.Error("Failed to update secret")
// 	}
// 	if secret.DateAdded == secret.DateModified {
// 		t.Error("Date modified not updated on secret update")
// 	}

// 	// Rename a secret
// 	secret, err = c.RenameSecret(secretList[1].name, "newname")
// 	if err != nil {
// 		t.Error("Failed to rename secret")
// 	} else {
// 		secretList[1].name = secret.Name
// 	}

// 	// Attempt renamed secret retrieval
// 	secret, err = c.GetSecret(secretList[1].name)
// 	if err != nil {
// 		t.Error("Secret rename failed to persist")
// 	}
// 	if secret.DateAdded == secret.DateModified {
// 		t.Error("Date modified not updated on secret rename")
// 	}

// 	// Rename a secret that doesn't exist
// 	secret, err = c.RenameSecret("invalidname", "newname")
// 	if err == nil {
// 		t.Error("Secret rename on non-existing secret did not fail")
// 	}

// 	// Rename to empty secret
// 	secret, err = c.RenameSecret("invalidname", "")
// 	if err == nil {
// 		t.Error("Secret rename with empty target did not fail")
// 	}

// 	// Test secret deletion
// 	// Middle
// 	testDeleteSecret(t, secretList[3].name, c)
// 	// Bottom
// 	testDeleteSecret(t, secretList[len(secretList)-1].name, c)
// 	// Top
// 	testDeleteSecret(t, secretList[0].name, c)

// 	// Secret does not exist
// 	testSecret = "invalidname"
// 	_, err = c.DeleteSecret(testSecret)
// 	if err == nil {
// 		t.Error("DeleteSecret on non-existing secret should return error", testSecret)
// 	}

// 	c.GetSecrets()

// 	c.SetFilename("")
// 	c.SetWriter(os.Stdout)
// 	c.Save()

// 	os.Remove(collectionFile)
// }

// ---

type errorReader struct{}

func (e errorReader) Read(p []byte) (int, error) {
	return 0, errors.New("error")
}

func TestNewCollection(t *testing.T) {
	tests := []struct {
		name string
		want *Collection
	}{
		{
			name: "new collection",
			want: &Collection{
				Secrets: make(map[string]Secret),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCollection(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCollectionWithData(t *testing.T) {
	addedTime, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	modifiedTime, _ := time.Parse(time.RFC3339, "2012-11-02T22:08:41+00:00")

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Collection
		wantErr bool
	}{
		{
			name: "new collection with data",
			args: args{
				data: []byte(`{ "Secrets": { "testname": { "DateAdded": "2012-11-01T22:08:41+00:00", "DateModified": "2012-11-02T22:08:41+00:00","Name": "testname", "Value": "seedseed" } } }`),
			},
			want: &Collection{
				Secrets: map[string]Secret{
					"testname": {
						DateAdded:    addedTime,
						DateModified: modifiedTime,
						Name:         "testname",
						Value:        "seedseed",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "fail to create new collection with data",
			args: args{
				data: []byte(`{`),
			},
			want: &Collection{
				Secrets: make(map[string]Secret),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCollectionWithData(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCollectionWithData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCollectionWithData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCollectionWithReader(t *testing.T) {
	addedTime, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	modifiedTime, _ := time.Parse(time.RFC3339, "2012-11-02T22:08:41+00:00")

	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *Collection
		wantErr bool
	}{
		{
			name: "new collection with reader",
			args: args{
				reader: strings.NewReader(`{ "Secrets": { "testname": { "DateAdded": "2012-11-01T22:08:41+00:00", "DateModified": "2012-11-02T22:08:41+00:00","Name": "testname", "Value": "seedseed" } } }`),
			},
			want: &Collection{
				Secrets: map[string]Secret{
					"testname": {
						DateAdded:    addedTime,
						DateModified: modifiedTime,
						Name:         "testname",
						Value:        "seedseed",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "fail new collection with reader with invalid data",
			args: args{
				reader: strings.NewReader(`{`),
			},
			want: &Collection{
				Secrets: make(map[string]Secret),
			},
			wantErr: true,
		},
		{
			name: "fail new collection with reader that returns errors",
			args: args{
				reader: errorReader{},
			},
			want: &Collection{
				Secrets: make(map[string]Secret),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCollectionWithReader(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCollectionWithReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCollectionWithReader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCollectionWithFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    *Collection
		wantErr bool
	}{
		{
			name: "collection file does not exist",
			args: args{
				filename: "nosuchfile.json",
			},
			want: &Collection{
				filename: "nosuchfile.json",
			},
			wantErr: true,
		},
		{
			name: "collection file exists",
			args: args{
				filename: "testcollection.json",
			},
			want: &Collection{
				filename: "testcollection.json",
			},
			wantErr: false,
		},
	}

	createTestData(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCollectionWithFile(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCollectionWithFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.filename != tt.want.filename {
				t.Errorf("NewCollectionWithFile() = %v, want %v", got, tt.want)
			}
		})
	}

	os.Remove("testcollection.json")
}

func TestCollection_GenerateCode(t *testing.T) {
	addedTime, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	modifiedTime, _ := time.Parse(time.RFC3339, "2012-11-02T22:08:41+00:00")

	type fields struct {
		Secrets  map[string]Secret
		filename string
		writer   io.Writer
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "generate code for secret",
			fields: fields{
				Secrets: map[string]Secret{
					"testname": {
						DateAdded:    addedTime,
						DateModified: modifiedTime,
						Name:         "testname",
						Value:        "seedseed",
					},
				},
			},
			args: args{
				name: "testname",
			},
			want:    6,
			wantErr: false,
		},
		{
			name: "generate code for secret that does not exist",
			fields: fields{
				Secrets: map[string]Secret{
					"testname": {
						DateAdded:    addedTime,
						DateModified: modifiedTime,
						Name:         "testname",
						Value:        "seedseed",
					},
				},
			},
			args: args{
				name: "invalidname",
			},
			want:    6,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collection{
				Secrets:  tt.fields.Secrets,
				filename: tt.fields.filename,
				writer:   tt.fields.writer,
			}
			got, err := c.GenerateCode(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Collection.GenerateCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if l := len(got); err == nil && l != 6 {
				t.Errorf("Collection.GenerateCode() length = %v, want %v", l, tt.want)
			}

			if _, atoiErr := strconv.Atoi(got); err == nil && atoiErr != nil {
				t.Errorf("Collection.GenerateCode() int conversion failed: %v", atoiErr)
			}
		})
	}
}

func TestCollection_GetSecrets(t *testing.T) {
	addedTime, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	modifiedTime, _ := time.Parse(time.RFC3339, "2012-11-02T22:08:41+00:00")

	type fields struct {
		Secrets  map[string]Secret
		filename string
		writer   io.Writer
	}
	tests := []struct {
		name   string
		fields fields
		want   []Secret
	}{
		{
			name: "get secrets",
			fields: fields{
				Secrets: map[string]Secret{
					"testname": {
						DateAdded:    addedTime,
						DateModified: modifiedTime,
						Name:         "testname",
						Value:        "seedseed",
					},
				},
			},
			want: []Secret{
				{
					DateAdded:    addedTime,
					DateModified: modifiedTime,
					Name:         "testname",
					Value:        "seedseed",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collection{
				Secrets:  tt.fields.Secrets,
				filename: tt.fields.filename,
				writer:   tt.fields.writer,
			}
			if got := c.GetSecrets(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collection.GetSecrets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollection_SetFilename(t *testing.T) {
	addedTime, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	modifiedTime, _ := time.Parse(time.RFC3339, "2012-11-02T22:08:41+00:00")

	type fields struct {
		Secrets  map[string]Secret
		filename string
		writer   io.Writer
	}
	type args struct {
		filename string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "set filename success",
			fields: fields{
				Secrets: map[string]Secret{
					"testname": {
						DateAdded:    addedTime,
						DateModified: modifiedTime,
						Name:         "testname",
						Value:        "seedseed",
					},
				},
			},
			args: args{
				filename: "testfile",
			},
			want: "testfile",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collection{
				Secrets:  tt.fields.Secrets,
				filename: tt.fields.filename,
				writer:   tt.fields.writer,
			}
			if got := c.SetFilename(tt.args.filename); got != tt.want {
				t.Errorf("Collection.SetFilename() = %v, want %v", got, tt.want)
			}
			if got := c.filename; got != tt.want {
				t.Errorf("Collection.SetFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollection_SetWriter(t *testing.T) {
	addedTime, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	modifiedTime, _ := time.Parse(time.RFC3339, "2012-11-02T22:08:41+00:00")

	type fields struct {
		Secrets  map[string]Secret
		filename string
		writer   io.Writer
	}
	tests := []struct {
		name       string
		fields     fields
		wantWriter string
	}{
		{
			name: "set filename success",
			fields: fields{
				Secrets: map[string]Secret{
					"testname": {
						DateAdded:    addedTime,
						DateModified: modifiedTime,
						Name:         "testname",
						Value:        "seedseed",
					},
				},
			},
			wantWriter: "testdata",
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collection{
				Secrets:  tt.fields.Secrets,
				filename: tt.fields.filename,
				writer:   tt.fields.writer,
			}
			writer := &bytes.Buffer{}
			writer.Write([]byte("testdata"))
			c.SetWriter(writer)
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("Collection.SetWriter() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestCollection_UpdateSecret(t *testing.T) {
	addedTime, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	modifiedTime, _ := time.Parse(time.RFC3339, "2012-11-02T22:08:41+00:00")

	type fields struct {
		Secrets  map[string]Secret
		filename string
		writer   io.Writer
	}
	type args struct {
		name  string
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Secret
		wantErr bool
	}{
		{
			name: "update (add) secret success",
			fields: fields{
				Secrets: map[string]Secret{
					"testname": {
						DateAdded:    addedTime,
						DateModified: modifiedTime,
						Name:         "testname",
						Value:        "seedseed",
					},
				},
			}, args: args{
				name:  "newname",
				value: "seed",
			},
			want: Secret{
				Name:  "newname",
				Value: "seed",
			},
			wantErr: false,
		},
		{
			name: "update existing secret success",
			fields: fields{
				Secrets: map[string]Secret{
					"testname": {
						DateAdded:    addedTime,
						DateModified: modifiedTime,
						Name:         "testname",
						Value:        "seedseed",
					},
				},
			}, args: args{
				name:  "testname",
				value: "seed",
			},
			want: Secret{
				Name:  "testname",
				Value: "seed",
			},
			wantErr: false,
		},
		{
			name: "update secret with empty name",
			fields: fields{
				Secrets: map[string]Secret{
					"testname": {
						DateAdded:    addedTime,
						DateModified: modifiedTime,
						Name:         "testname",
						Value:        "seedseed",
					},
				},
			}, args: args{
				name:  "",
				value: "seed",
			},
			want:    Secret{},
			wantErr: true,
		},
		{
			name: "update secret with empty value",
			fields: fields{
				Secrets: map[string]Secret{
					"testname": {
						DateAdded:    addedTime,
						DateModified: modifiedTime,
						Name:         "testname",
						Value:        "seedseed",
					},
				},
			}, args: args{
				name:  "newname",
				value: "",
			},
			want:    Secret{},
			wantErr: true,
		},
		{
			name: "update secret with invalid value",
			fields: fields{
				Secrets: map[string]Secret{
					"testname": {
						DateAdded:    addedTime,
						DateModified: modifiedTime,
						Name:         "testname",
						Value:        "seedseed",
					},
				},
			}, args: args{
				name:  "newname",
				value: "#$%^&*(",
			},
			want:    Secret{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collection{
				Secrets:  tt.fields.Secrets,
				filename: tt.fields.filename,
				writer:   tt.fields.writer,
			}
			got, err := c.UpdateSecret(tt.args.name, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Collection.UpdateSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, err = c.GetSecret(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Collection.UpdateSecret() with Collection.GetSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Name != tt.want.Name || got.Value != tt.want.Value {
				t.Errorf("Collection.UpdateSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollection_RenameSecret(t *testing.T) {
	addedTime, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	modifiedTime, _ := time.Parse(time.RFC3339, "2012-11-02T22:08:41+00:00")

	type fields struct {
		Secrets  map[string]Secret
		filename string
		writer   io.Writer
	}
	type args struct {
		oldName string
		newName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Secret
		wantErr bool
	}{
		{
			name: "rename secret success",
			fields: fields{
				Secrets: map[string]Secret{
					"testname": {
						DateAdded:    addedTime,
						DateModified: modifiedTime,
						Name:         "testname",
						Value:        "seedseed",
					},
				},
			}, args: args{
				oldName: "testname",
				newName: "newname",
			},
			want: Secret{
				Name:  "newname",
				Value: "seedseed",
			},
			wantErr: false,
		},
		{
			name: "rename secret new name empty",
			fields: fields{
				Secrets: map[string]Secret{
					"testname": {
						DateAdded:    addedTime,
						DateModified: modifiedTime,
						Name:         "testname",
						Value:        "seedseed",
					},
				},
			}, args: args{
				oldName: "testname",
				newName: "",
			},
			want:    Secret{},
			wantErr: true,
		},
		{
			name: "rename secret old name not found",
			fields: fields{
				Secrets: map[string]Secret{
					"testname": {
						DateAdded:    addedTime,
						DateModified: modifiedTime,
						Name:         "testname",
						Value:        "seedseed",
					},
				},
			}, args: args{
				oldName: "invalidname",
				newName: "newname",
			},
			want:    Secret{},
			wantErr: true,
		},
		{
			name: "rename secret old name empty",
			fields: fields{
				Secrets: map[string]Secret{
					"testname": {
						DateAdded:    addedTime,
						DateModified: modifiedTime,
						Name:         "testname",
						Value:        "seedseed",
					},
				},
			}, args: args{
				oldName: "",
				newName: "newname",
			},
			want:    Secret{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collection{
				Secrets:  tt.fields.Secrets,
				filename: tt.fields.filename,
				writer:   tt.fields.writer,
			}
			got, err := c.RenameSecret(tt.args.oldName, tt.args.newName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Collection.RenameSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Name != tt.want.Name || got.Value != tt.want.Value {
				t.Errorf("Collection.RenameSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollection_DeleteSecret(t *testing.T) {
	addedTime, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	modifiedTime, _ := time.Parse(time.RFC3339, "2012-11-02T22:08:41+00:00")

	type fields struct {
		Secrets  map[string]Secret
		filename string
		writer   io.Writer
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Secret
		wantErr bool
	}{
		{
			name: "delete secret not found",
			fields: fields{
				Secrets: map[string]Secret{
					"testname": {
						DateAdded:    addedTime,
						DateModified: modifiedTime,
						Name:         "testname",
						Value:        "seedseed",
					},
				},
			}, args: args{
				name: "invalidname",
			},
			want:    Secret{},
			wantErr: true,
		},
		{
			name: "delete secret success",
			fields: fields{
				Secrets: map[string]Secret{
					"testname": {
						DateAdded:    addedTime,
						DateModified: modifiedTime,
						Name:         "testname",
						Value:        "seedseed",
					},
				},
			}, args: args{
				name: "testname",
			},
			want: Secret{
				Name:  "testname",
				Value: "seedseed",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collection{
				Secrets:  tt.fields.Secrets,
				filename: tt.fields.filename,
				writer:   tt.fields.writer,
			}
			got, err := c.DeleteSecret(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Collection.DeleteSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Name != tt.want.Name || got.Value != tt.want.Value {
				t.Errorf("Collection.DeleteSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollection_Save(t *testing.T) {
	addedTime, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	modifiedTime, _ := time.Parse(time.RFC3339, "2012-11-02T22:08:41+00:00")

	type fields struct {
		Secrets  map[string]Secret
		filename string
		writer   io.Writer
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "save collection to filename success",
			fields: fields{
				Secrets: map[string]Secret{
					"testname": {
						DateAdded:    addedTime,
						DateModified: modifiedTime,
						Name:         "testname",
						Value:        "seedseed",
					},
				},
				filename: "testcollection.json",
			},
			wantErr: false,
		},
		{
			name: "save collection to writer success",
			fields: fields{
				Secrets: map[string]Secret{
					"testname": {
						DateAdded:    addedTime,
						DateModified: modifiedTime,
						Name:         "testname",
						Value:        "seedseed",
					},
				},
				writer: &bytes.Buffer{},
			},
			wantErr: false,
		},
		{
			name: "save collection no filename or writer failure",
			fields: fields{
				Secrets: map[string]Secret{
					"testname": {
						DateAdded:    addedTime,
						DateModified: modifiedTime,
						Name:         "testname",
						Value:        "seedseed",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collection{
				Secrets:  tt.fields.Secrets,
				filename: tt.fields.filename,
				writer:   tt.fields.writer,
			}
			if err := c.Save(); (err != nil) != tt.wantErr {
				t.Errorf("Collection.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	os.Remove("testcollection.json")
}
