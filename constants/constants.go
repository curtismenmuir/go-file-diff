package constants

// Error messages
const (
	ModeFlagMissingError               string = "Error: Must set at least one mode"
	SignatureFlagsMissingError         string = "Error: Must provide Original & Signature files when enabling Signature mode"
	DeltaFlagsMissingError             string = "Error: Must provide Signature, Updated & Delta files when enabling Delta mode"
	SignatureDeltaFlagsMissingError    string = "Error: Must provide Updated & Delta files when enabling Signature & Delta modes"
	DeltaNotImplementedError           string = "Error: Delta mode not implemented, coming soon"
	UnableToCheckFileFolderExistsError string = "Error: Unable to check if file/folder exists"
	FileDoesNotExistError              string = "Error: File does not exist"
	OriginalFileDoesNotExistError      string = "Error: Original file does not exist"
	SearchingForFileButFoundDirError   string = "Error: Searching for a file but found a folder dir"
	OriginalFileIsFolderError          string = "Error: Original file provided is a folder dir"
	UnableToCreateNewFolderError       string = "Error: Unable to create folder"
	UnableToCreateOutputsFolder        string = "Error: Unable to create Outputs folder"
	UnableToCreateSignatureFile        string = "Error: Unable to create Signature file"
	UnableToWriteToSignatureFile       string = "Error: Unable to write to Signature file"
	EndOfFileError                     string = "Error: Reached EOF"
	UnableToGenerateSignature          string = "Error: Unable to generate Signature"
	SignatureFileDoesNotExistError     string = "Error: Signature file does not exist"
	UnableToOpenSignatureFile          string = "Error: Unable to open Signature file"
	UnableToDecodeSignatureFromFile    string = "Error: Unable to decode Signature from file"
)
