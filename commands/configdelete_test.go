package commands

import (
	"os"
	"testing"

	"github.com/arcanericky/totp"
)

func TestConfigDelete(t *testing.T) {
	defaults()
	collectionFile.filename = "testcollection.json"

	secretList := createTestData(t)

	configDeleteCmd := getConfigDeleteCmd()
	configDeleteCmd.Run(nil, []string{"secret"})

	_ = configDeleteCmd.Flags().Set(optionYes, "true")

	// Secret does not exit
	configDeleteCmd.Run(nil, []string{"secret"})

	// No secret provided
	configDeleteCmd.Run(nil, []string{})

	// Successful delete
	configDeleteCmd.Run(nil, []string{secretList[3].name})
	c, err := totp.NewCollectionWithFile(collectionFile.filename)
	if err != nil {
		t.Error("Could not load collection for delete test from file")
	}

	_, err = c.GetSecret(secretList[3].name)
	if err == nil {
		t.Error("Secret not deleted")
	}

	// No collections file
	os.Remove(collectionFile.filename)
	deleteSecret(secretList[3].name)
}
