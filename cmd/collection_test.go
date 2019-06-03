package cmd

import (
	"testing"

	"github.com/arcanericky/totp"
)

type secretItem struct {
	name  string
	value string
}

func createTestData(t *testing.T) []secretItem {
	t.Helper()

	// Create test data
	c, _ := totp.NewCollectionWithFile(defaultCollectionFile)

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

	c.Save()

	return secretList
}
