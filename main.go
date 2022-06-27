package main

import (
	"errors"
	"fmt"

	"github.com/curtismenmuir/go-file-diff/cmd"
	"github.com/curtismenmuir/go-file-diff/constants"
	"github.com/curtismenmuir/go-file-diff/utils"
)

var (
	logger    = utils.Logger
	parseCMD  = cmd.ParseCMD
	verifyCMD = cmd.VerifyCMD
)

func getSignature(verbose bool, originalFile string, signatureFile string) error {
	return errors.New(constants.SignatureNotImplementedError)
}

func getDelta(verbose bool, signatureFile string, updatedFile string, deltaFile string) error {
	return errors.New(constants.DeltaNotImplementedError)
}

func main() {
	// Parse CMD flags
	verbose, signatureMode, deltaMode, originalFile, signatureFile, updatedFile, deltaFile := parseCMD()
	// Verify valid CMD flags provided
	if !verifyCMD(verbose, signatureMode, deltaMode, originalFile, signatureFile, updatedFile, deltaFile) {
		return
	}

	// Generate Signature
	if signatureMode {
		err := getSignature(verbose, originalFile, signatureFile)
		if err != nil {
			logger(fmt.Sprintf("Error: %s", err.Error()), true)
			return
		}
	}

	// Generate Delta
	if deltaMode {
		err := getDelta(verbose, signatureFile, updatedFile, deltaFile)
		if err != nil {
			logger(fmt.Sprintf("Error: %s", err.Error()), true)
			return
		}
	}
}
