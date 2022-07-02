package commands

import (
	"testing"
)

func TestPrintResult(t *testing.T) {
	text := "test text"
	_, _ = printResultf(text)
	collectionFile.useStdio = true
	_, _ = printResultf(text)
}
