package totp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/pquerna/otp/totp"
)

// Entry is a struct containing data necessary for working with keys,
// namely, the name of the key and the seed value
type Entry struct {
	DateAdded    time.Time
	DateModified time.Time
	Name         string
	Seed         string
}

// Collection is a struct that holds settings data
type Collection struct {
	seeds    map[string]Entry
	filename string
}

// SettingsInterface is used for DI when needed
type SettingsInterface interface {
	Save() error

	DeleteKey(string) (Entry, error)
	GetKey(string) (Entry, error)
	GetKeys() []Entry
	SetFilename(string) string
	UpdateKey(string, string) (Entry, error)
}

// Save serializes (marshals) the entire TotpSettings struct and writes it to
// a file
func (c *Collection) Save() error {
	var err error

	if len(c.filename) == 0 {
		err = errors.New("no filename configured")
	} else {
		serializedSettings, err := c.Serialize()
		if err == nil {
			err = ioutil.WriteFile(c.filename, serializedSettings, 0600)
		}
	}

	return err
}

// DeleteKey deletes an Entry by name
func (c *Collection) DeleteKey(name string) (Entry, error) {
	var err error

	retKey, ok := c.seeds[name]

	if ok == true {
		delete(c.seeds, name)
	} else {
		err = errors.New("Entry does not exist")
	}

	return retKey, err
}

// UpdateKey updates (if it exists) or adds a new Entry with the
// name and seed given
func (c *Collection) UpdateKey(name, seed string) (Entry, error) {
	var retKey Entry
	var err error

	if len(name) == 0 {
		err = errors.New("Key name must not be empty")
	} else if len(seed) == 0 {
		err = errors.New("Key seed must not be empty")
	} else {
		_, err = totp.GenerateCode(seed, time.Now())
		if err == nil {
			retKey, ok := c.seeds[name]
			if ok == true {
				retKey.Seed = seed
				retKey.DateModified = time.Now()
				c.seeds[name] = retKey
			} else {
				dateAdded := time.Now()
				newKey := Entry{Name: name, Seed: seed, DateAdded: dateAdded, DateModified: dateAdded}
				c.seeds[name] = newKey
				retKey = newKey
			}
		}
	}

	return retKey, err
}

// RenameKey renames a key
func (c *Collection) RenameKey(oldName, newName string) (Entry, error) {
	var retKey Entry
	var ok bool
	var err error

	if len(newName) != 0 {
		retKey, ok = c.seeds[oldName]
		if ok == true {
			retKey.Name = newName
			retKey.DateModified = time.Now()
			c.seeds[newName] = retKey
			delete(c.seeds, oldName)
		} else {
			err = errors.New("Key not found")
		}
	} else {
		err = errors.New("Key name must not be empty")
	}

	return retKey, err
}

// GetKey returns an Entry with the name argument
func (c *Collection) GetKey(name string) (Entry, error) {
	var err error

	retKey, ok := c.seeds[name]
	if ok == false {
		err = fmt.Errorf("Key name \"%s\" not found", name)
	}

	return retKey, err
}

// GetKeys returns a slice containing all the keys
func (c *Collection) GetKeys() []Entry {
	keys := []Entry{}
	for _, key := range c.seeds {
		keys = append(keys, key)
	}

	return keys
}

// GenerateCodeWithTime creates a TOTP code with the named entry's seed
func (c *Collection) GenerateCodeWithTime(name string, time time.Time) (string, error) {
	var code string

	key, err := c.GetKey(name)
	if err == nil {
		code, err = totp.GenerateCode(key.Seed, time)
	}

	return code, err
}

// GenerateCode creates a TOTP code with the named entry's seed
func (c *Collection) GenerateCode(name string) (string, error) {
	return c.GenerateCodeWithTime(name, time.Now())
}

// Serialize marshals a Entries struct into a byte array
func (c *Collection) Serialize() ([]byte, error) {
	return json.MarshalIndent(c.seeds, "", "  ")
}

// Deserialize unmarshals a byte array into a Entries struct
func (c *Collection) Deserialize(data []byte) error {
	return json.Unmarshal(data, &c.seeds)
}

// SetFilename sets the filename for the Save method
func (c *Collection) SetFilename(filename string) string {
	c.filename = filename

	return c.filename
}

// NewCollection creates a new, blank Collection instance
func NewCollection() *Collection {
	c := new(Collection)
	c.seeds = make(map[string]Entry)
	return c
}

// NewCollectionWithData creates a new Collection instance with data from a byte slice
func NewCollectionWithData(data []byte) (*Collection, error) {
	c := NewCollection()
	err := c.Deserialize(data)

	return c, err
}

// NewCollectionWithFile creates a new Collection instance with data from a file
func NewCollectionWithFile(filename string) (*Collection, error) {
	data, err := ioutil.ReadFile(filename)

	c := NewCollection()

	if err == nil {
		c, err = NewCollectionWithData(data)
	}

	c.filename = filename

	return c, err
}
