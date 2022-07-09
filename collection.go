package totp

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"

	"github.com/pquerna/otp/totp"
)

var ErrSecretNotFound = errors.New("secret not found")
var ErrNoFilename = errors.New("no save target")
var ErrSecretNameEmpty = errors.New("secret name empty")
var ErrSecretValueEmpty = errors.New("secret value empty")

// Secret is a struct containing data necessary for working with secrets,
// namely, the name of the secret name and the secret value
type Secret struct {
	// DateAdded is the date a secret was added to the collection
	DateAdded time.Time

	// DateModified is the date a secret was last modified
	DateModified time.Time

	// Name is the name of the secret used for retrieval
	Name string

	// Value is the secret (seed) value
	Value string
}

// Collection is a struct that holds TOTP data
type Collection struct {
	// Secrets is a map of secrets using the secret name as the key
	Secrets map[string]Secret

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
	if err != nil {
		return err
	}

	if c.writer != nil {
		_, err = c.writer.Write(serializedSettings)
	} else if len(c.filename) != 0 {
		err = os.WriteFile(c.filename, serializedSettings, 0600)
	} else {
		err = ErrNoFilename
	}

	return err
}

// DeleteSecret deletes an entry by name
func (c *Collection) DeleteSecret(name string) (Secret, error) {
	retSecret, ok := c.Secrets[name]
	if !ok {
		return Secret{}, ErrSecretNotFound
	}

	delete(c.Secrets, name)

	return retSecret, nil
}

// UpdateSecret updates (if it exists) or adds a new entry with the
// name and value given
func (c *Collection) UpdateSecret(name, value string) (Secret, error) {
	if len(name) == 0 {
		return Secret{}, ErrSecretNameEmpty
	}

	if len(value) == 0 {
		return Secret{}, ErrSecretValueEmpty
	}

	_, err := totp.GenerateCode(value, time.Now())
	if err != nil {
		return Secret{}, err
	}

	retSecret, ok := c.Secrets[name]
	if ok {
		// entry indicates an update
		retSecret.Value = value
		retSecret.DateModified = time.Now()
		c.Secrets[name] = retSecret
	} else {
		// no entry indicates an add
		dateAdded := time.Now()
		retSecret = Secret{
			Name:         name,
			Value:        value,
			DateAdded:    dateAdded,
			DateModified: dateAdded,
		}
		c.Secrets[name] = retSecret
	}

	return retSecret, err
}

// RenameSecret renames a secret
func (c *Collection) RenameSecret(oldName, newName string) (Secret, error) {
	if len(newName) == 0 {
		return Secret{}, ErrSecretNameEmpty
	}

	retSecret, ok := c.Secrets[oldName]
	if !ok {
		return Secret{}, ErrSecretNotFound
	}

	retSecret.Name = newName
	retSecret.DateModified = time.Now()
	c.Secrets[newName] = retSecret
	delete(c.Secrets, oldName)

	return retSecret, nil
}

// GetSecret returns a secret with the name argument
func (c *Collection) GetSecret(name string) (Secret, error) {
	retSecret, ok := c.Secrets[name]
	if !ok {
		return Secret{}, ErrSecretNotFound
	}

	return retSecret, nil
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
	secret, err := c.GetSecret(name)
	if err != nil {
		return "", err
	}

	return totp.GenerateCode(secret.Value, time)

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

	return c, c.Deserialize(data)
}

// NewCollectionWithReader creates a new collection from a Reader interface
func NewCollectionWithReader(reader io.Reader) (*Collection, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return NewCollection(), err
	}

	return NewCollectionWithData(data)
}

// NewCollectionWithFile creates a new Collection instance with data from a file.
// If the file open fails, a new Collection instance is returned along with the
// file open error, which guarantees a usable but empty collection is returned.
//
// Returning data to be used along with an error is bad design but changing this
// would be a breaking API change.
func NewCollectionWithFile(filename string) (c *Collection, err error) {
	f, err := os.Open(filename)
	if err != nil {
		c := NewCollection()
		c.filename = filename
		return c, err
	}

	defer func() {
		if e := f.Close(); e != nil && err == nil {
			err = e
		}
	}()

	c, err = NewCollectionWithReader(f)
	c.filename = filename

	return c, err
}
