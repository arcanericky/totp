package commands

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
	c, _ := totp.NewCollectionWithFile(collectionFile.filename)

	// Create some test data
	secretList := []secretItem{
		{name: "name0", value: "SEED"},
		{name: "name1", value: "SEED"},
		{name: "name2", value: "SEEDSEED"},
		{name: "name3", value: "SEED"},
		{name: "name4", value: "SEED"},
		{name: "testname", value: "TESTSECRET"},
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
