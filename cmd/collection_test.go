package cmd

import (
	"testing"

	"github.com/arcanericky/totp"
)

type seedItem struct {
	name string
	seed string
}

func createTestData(t *testing.T) []seedItem {
	t.Helper()

	// Create test data
	c, _ := totp.NewCollectionWithFile(defaultCollectionFile)

	// Create some test data
	seedList := []seedItem{
		{name: "name0", seed: "seed"},
		{name: "name1", seed: "seed"},
		{name: "name2", seed: "seedseed"},
		{name: "name3", seed: "seed"},
		{name: "name4", seed: "seed"},
	}

	for _, i := range seedList {
		_, err := c.UpdateKey(i.name, i.seed)
		if err != nil {
			t.Error("Error updating key:", err)
		}
	}

	c.Save()

	return seedList
}
