package cmd

import (
	"testing"
)

func TestPrintResult(t *testing.T) {
	text := "test text"
	printResultf(text)
	collectionFile.useStdio = true
	printResultf(text)
}
