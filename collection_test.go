package totp

import (
	"os"
	"testing"
)

func testDeleteKey(t *testing.T, testKey string, c *Collection) {
	t.Helper()

	_, err := c.DeleteKey(testKey)
	if err != nil {
		t.Error("DeleteKey on valid key returned error", err)
	}

	_, err = c.GetKey(testKey)
	if err == nil {
		t.Error("Error deleting key", testKey)
	}
}

func TestWriteProtected(t *testing.T) {
	name := "name"
	seed := "seed"

	updatedValue := "updatedvalue"

	s := NewCollection()

	key, err := s.UpdateKey(name, seed)
	if err != nil {
		t.Error("Failed to add key for test")
	}

	key.Name = updatedValue
	key.Seed = updatedValue

	// Validate key returned on update can't be changed
	key, err = s.GetKey(name)
	if key.Name != name || key.Seed != seed {
		t.Error("Internal settings key can be updated with returned key from UpdateKey()")
	}

	key.Name = updatedValue
	key.Seed = updatedValue

	// Validate key returned on get can't get changed
	key, err = s.GetKey(name)
	if key.Name != name || key.Seed != seed {
		t.Error("Internal settings key can be updated with returned key from UpdateKey()")
	}
}

func TestSettingsNew(t *testing.T) {
	collectionFile := "testcollection.json"

	c := NewCollection()

	// Test error on Save with no filename
	err := c.Save()
	if err == nil {
		t.Error("Save collection with no filename should generate error")
	}

	// Set filename for remainder of tests
	c.SetFilename(collectionFile)

	// Create some data
	type seedItem struct {
		name string
		seed string
	}

	// Create some test data
	seedList := []seedItem{
		{name: "name0", seed: "seed"},
		{name: "name1", seed: "seed"},
		{name: "name2", seed: "seed"},
		{name: "name3", seed: "seed"},
		{name: "name4", seed: "seed"},
	}

	for _, i := range seedList {
		_, err := c.UpdateKey(i.name, i.seed)
		if err != nil {
			t.Error("Error updating key:", err)
		}
	}

	err = c.Save()
	if err != nil {
		t.Error("Save collection with filename yielded error")
	}

	// Load test data
	c, _ = NewCollectionWithFile(collectionFile)
	for _, i := range seedList {
		key, err := c.GetKey(i.name)
		if err == nil {
			if key.Name != i.name || key.Seed != i.seed {
				t.Error("Loaded keys don't match saved keys")
			}
		} else {
			t.Error("Error loading test data:", err)
		}
	}

	// Test GenerateCode() methods
	testKey := seedList[0].name
	_, err = c.GenerateCode(testKey)
	if err != nil {
		t.Error("Error generating code for key", testKey)
	}

	// Attempt invalid key retrieval
	key, err := c.GetKey("invalidkey")
	if err == nil {
		t.Error("GetKey returned success on invalid key retrieval")
	}

	newSeed := "deadbeef"

	// Update key with empty name
	key, err = c.UpdateKey("", newSeed)
	if err == nil {
		t.Error("Update key with empty name did not return error")
	}

	// Update key with empty seed
	key, err = c.UpdateKey(seedList[0].name, "")
	if err == nil {
		t.Error("Update key with empty seed did not return error")
	}

	// Update a key
	testKey = seedList[0].name
	key, err = c.UpdateKey(testKey, newSeed)
	if err != nil {
		t.Error("Error updating key", key, err)
	}
	key, err = c.GetKey(testKey)
	if err != nil || key.Seed != newSeed {
		t.Error("Failed to update key")
	}
	if key.DateAdded == key.DateModified {
		t.Error("Date modified not updated on key update")
	}

	// Rename a key
	key, err = c.RenameKey(seedList[1].name, "newname")
	if err != nil {
		t.Error("Failed to rename key")
	} else {
		seedList[1].name = key.Name
	}

	// Attempt renamed key retrieval
	key, err = c.GetKey(seedList[1].name)
	if err != nil {
		t.Error("Key rename failed to persist")
	}
	if key.DateAdded == key.DateModified {
		t.Error("Date modified not updated on key rename")
	}

	// Rename a key that doesn't exist
	key, err = c.RenameKey("invalidname", "newname")
	if err == nil {
		t.Error("Key rename on non-existing key did not fail")
	}

	// Rename to empty key
	key, err = c.RenameKey("invalidname", "")
	if err == nil {
		t.Error("Key rename with empty target did not fail")
	}

	// Test key deletion
	// Middle
	testDeleteKey(t, seedList[3].name, c)
	// Bottom
	testDeleteKey(t, seedList[len(seedList)-1].name, c)
	// Top
	testDeleteKey(t, seedList[0].name, c)

	// Key does not exist
	testKey = "invalidkey"
	_, err = c.DeleteKey(testKey)
	if err == nil {
		t.Error("DeleteKey on non-existing key should return error", testKey)
	}

	c.GetKeys()

	os.Remove(collectionFile)
}
