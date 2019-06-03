package totp

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/pquerna/otp/totp"
)

var errSecretNotFound = errors.New("Secret not found")
var errNoFilename = errors.New("No filename")
var errSecretNameEmpty = errors.New("Secret name empty")
var errSecretValueEmpty = errors.New("Secret value empty")

// Secret is a struct containing data necessary for working with secrets,
// namely, the name of the secret name and the secret value
type Secret struct {
	DateAdded    time.Time
	DateModified time.Time
	Name         string
	Value        string
}

// Collection is a struct that holds TOTP data
type Collection struct {
	Secrets  map[string]Secret
	filename string
	writer   io.Writer
}

// CollectionInterface is used for DI when needed
type CollectionInterface interface {
	DeleteSecret(string) (Secret, error)
	GetSecret(string) (Secret, error)
	GetSecrets() []Secret
	Save() error
	SetFilename(string) string
	UpdateSecret(string, string) (Secret, error)
}

// Save serializes (marshals) the Collections struct and writes it to
// a file
func (c *Collection) Save() error {
	serializedSettings, err := c.Serialize()

	if err == nil {
		if len(c.filename) == 0 && c.writer != nil {
			c.writer.Write(serializedSettings)
		} else {
			err = ioutil.WriteFile(c.filename, serializedSettings, 0600)
		}
	}

	return err
}

// DeleteSecret deletes an Entry by name
func (c *Collection) DeleteSecret(name string) (Secret, error) {
	var err error

	retSecret, ok := c.Secrets[name]

	if ok == true {
		delete(c.Secrets, name)
	} else {
		err = errSecretNotFound
	}

	return retSecret, err
}

// UpdateSecret updates (if it exists) or adds a new Entry with the
// name and value given
func (c *Collection) UpdateSecret(name, value string) (Secret, error) {
	var retSecret Secret
	var err error
	var ok bool

	if len(name) == 0 {
		err = errSecretNameEmpty
	} else if len(value) == 0 {
		err = errSecretValueEmpty
	} else {
		_, err = totp.GenerateCode(value, time.Now())
		if err == nil {
			retSecret, ok = c.Secrets[name]
			if ok == true {
				retSecret.Value = value
				retSecret.DateModified = time.Now()
				c.Secrets[name] = retSecret
			} else {
				dateAdded := time.Now()
				newSecret := Secret{Name: name, Value: value, DateAdded: dateAdded, DateModified: dateAdded}
				c.Secrets[name] = newSecret
				retSecret = newSecret
			}
		}
	}

	return retSecret, err
}

// RenameSecret renames a secret
func (c *Collection) RenameSecret(oldName, newName string) (Secret, error) {
	var retSecret Secret
	var ok bool
	var err error

	if len(newName) != 0 {
		retSecret, ok = c.Secrets[oldName]
		if ok == true {
			retSecret.Name = newName
			retSecret.DateModified = time.Now()
			c.Secrets[newName] = retSecret
			delete(c.Secrets, oldName)
		} else {
			err = errSecretNotFound
		}
	} else {
		err = errSecretNameEmpty
	}

	return retSecret, err
}

// GetSecret returns an Secret with the name argument
func (c *Collection) GetSecret(name string) (Secret, error) {
	var err error

	retSecret, ok := c.Secrets[name]
	if ok == false {
		err = errSecretNotFound
	}

	return retSecret, err
}

// GetSecrets returns a slice containing all the secrets
func (c *Collection) GetSecrets() []Secret {
	secrets := []Secret{}
	for _, secret := range c.Secrets {
		secrets = append(secrets, secret)
	}

	return secrets
}

// GenerateCodeWithTime creates a TOTP code with the named secret's value
func (c *Collection) GenerateCodeWithTime(name string, time time.Time) (string, error) {
	var code string

	secret, err := c.GetSecret(name)
	if err == nil {
		code, err = totp.GenerateCode(secret.Value, time)
	}

	return code, err
}

// GenerateCode creates a TOTP code with the named secret's value
func (c *Collection) GenerateCode(name string) (string, error) {
	return c.GenerateCodeWithTime(name, time.Now())
}

// Serialize marshals the Collection struct into a byte array
func (c *Collection) Serialize() ([]byte, error) {
	return json.MarshalIndent(c, "", "  ")
}

// Deserialize unmarshals a byte array into a Collection struct
func (c *Collection) Deserialize(data []byte) error {
	return json.Unmarshal(data, &c)
}

// SetWriter sets the writer for the Save method
func (c *Collection) SetWriter(writer io.Writer) {
	c.writer = writer
}

// SetFilename sets the filename for the Save method
func (c *Collection) SetFilename(filename string) string {
	c.filename = filename

	return c.filename
}

// NewCollection creates a new, blank Collection instance
func NewCollection() *Collection {
	c := new(Collection)
	c.Secrets = make(map[string]Secret)
	return c
}

// NewCollectionWithData creates a new Collection instance with data from a byte slice
func NewCollectionWithData(data []byte) (*Collection, error) {
	c := NewCollection()
	err := c.Deserialize(data)

	return c, err
}

// NewCollectionWithReader creates a new collection from a Reader interface
func NewCollectionWithReader(reader io.Reader) (*Collection, error) {
	data, err := ioutil.ReadAll(reader)

	if err != nil {
		return NewCollection(), err
	}

	return NewCollectionWithData(data)
}

// NewCollectionWithFile creates a new Collection instance with data from a file
func NewCollectionWithFile(filename string) (*Collection, error) {
	c := NewCollection()

	f, err := os.Open(filename)

	if err == nil {
		c, err = NewCollectionWithReader(f)
	}

	c.filename = filename

	return c, err
}
