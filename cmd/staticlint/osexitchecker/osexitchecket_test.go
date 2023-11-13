package osexitchecker

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestOsExitChecker(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), OsExitChecker, "./...")
}
