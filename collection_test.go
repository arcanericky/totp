package totp

import (
	"os"
	"testing"
)

func testDeleteSecret(t *testing.T, testSecret string, c *Collection) {
	t.Helper()

	_, err := c.DeleteSecret(testSecret)
	if err != nil {
		t.Error("DeleteSecret on valid secret name returned error", err)
	}

	_, err = c.GetSecret(testSecret)
	if err == nil {
		t.Error("Error deleting secret", testSecret)
	}
}

func TestWriteProtected(t *testing.T) {
	secretName := "name"
	secretValue := "seed"

	updatedValue := "updatedvalue"

	s := NewCollection()

	secret, err := s.UpdateSecret(secretName, secretValue)
	if err != nil {
		t.Error("Failed to add secret for test")
	}

	secret.Name = updatedValue
	secret.Value = updatedValue

	// Validate secret returned on update can't be changed
	secret, err = s.GetSecret(secretName)
	if secret.Name != secretName || secret.Value != secretValue {
		t.Error("Internal collection secret can be updated with returned secret from UpdateSecret()")
	}
}

func TestSettingsNew(t *testing.T) {
	collectionFile := "testcollection.json"

	// Test failure on Reader interface
	c, err := NewCollectionWithReader(os.Stdout)
	if err == nil {
		t.Error("New collection should fail with os.Stdout as reader")
	}

	c = NewCollection()

	// Test error on Save with no filename
	err = c.Save()
	if err == nil {
		t.Error("Save collection with no filename should generate error")
	}

	// Set filename for remainder of tests
	c.SetFilename(collectionFile)

	// Create some data
	type secretItem struct {
		name  string
		value string
	}

	// Create some test data
	secretList := []secretItem{
		{name: "name0", value: "seed"},
		{name: "name1", value: "seed"},
		{name: "name2", value: "seed"},
		{name: "name3", value: "seed"},
		{name: "name4", value: "seed"},
	}

	for _, i := range secretList {
		_, err := c.UpdateSecret(i.name, i.value)
		if err != nil {
			t.Error("Error updating secret:", err)
		}
	}

	err = c.Save()
	if err != nil {
		t.Error("Save collection with filename yielded error")
	}

	// Load test data
	c, _ = NewCollectionWithFile(collectionFile)
	for _, i := range secretList {
		secret, err := c.GetSecret(i.name)
		if err == nil {
			if secret.Name != i.name || secret.Value != i.value {
				t.Error("Loaded secrets don't match saved secrets")
			}
		} else {
			t.Error("Error loading test data:", err)
		}
	}

	// Test GenerateCode() methods
	testSecret := secretList[0].name
	_, err = c.GenerateCode(testSecret)
	if err != nil {
		t.Error("Error generating code for secret", testSecret)
	}

	// Attempt invalid secret retrieval
	secret, err := c.GetSecret("invalidsecret")
	if err == nil {
		t.Error("GetSecret returned success on invalid secret retrieval")
	}

	newSecret := "deadbeef"

	// Update secret with empty name
	secret, err = c.UpdateSecret("", newSecret)
	if err == nil {
		t.Error("UpdateSecret with empty name did not return error")
	}

	// Update secret with empty value
	secret, err = c.UpdateSecret(secretList[0].name, "")
	if err == nil {
		t.Error("UpdateSecret with empty value did not return error")
	}

	// Update a secret
	testSecret = secretList[0].name
	secret, err = c.UpdateSecret(testSecret, newSecret)
	if err != nil {
		t.Error("Error updating secret", secret, err)
	}
	if secret.DateAdded == secret.DateModified {
		t.Error("Date modified not updated on secret update")
	}

	secret, err = c.GetSecret(testSecret)
	if err != nil || secret.Value != newSecret {
		t.Error("Failed to update secret")
	}
	if secret.DateAdded == secret.DateModified {
		t.Error("Date modified not updated on secret update")
	}

	// Rename a secret
	secret, err = c.RenameSecret(secretList[1].name, "newname")
	if err != nil {
		t.Error("Failed to rename secret")
	} else {
		secretList[1].name = secret.Name
	}

	// Attempt renamed secret retrieval
	secret, err = c.GetSecret(secretList[1].name)
	if err != nil {
		t.Error("Secret rename failed to persist")
	}
	if secret.DateAdded == secret.DateModified {
		t.Error("Date modified not updated on secret rename")
	}

	// Rename a secret that doesn't exist
	secret, err = c.RenameSecret("invalidname", "newname")
	if err == nil {
		t.Error("Secret rename on non-existing secret did not fail")
	}

	// Rename to empty secret
	secret, err = c.RenameSecret("invalidname", "")
	if err == nil {
		t.Error("Secret rename with empty target did not fail")
	}

	// Test secret deletion
	// Middle
	testDeleteSecret(t, secretList[3].name, c)
	// Bottom
	testDeleteSecret(t, secretList[len(secretList)-1].name, c)
	// Top
	testDeleteSecret(t, secretList[0].name, c)

	// Secret does not exist
	testSecret = "invalidname"
	_, err = c.DeleteSecret(testSecret)
	if err == nil {
		t.Error("DeleteSecret on non-existing secret should return error", testSecret)
	}

	c.GetSecrets()

	c.SetFilename("")
	c.SetWriter(os.Stdout)
	c.Save()

	os.Remove(collectionFile)
}
