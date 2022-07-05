package constants

// Error messages
const (
	ModeFlagMissingError               string = "Error: Must set at least one mode"
	SignatureFlagsMissingError         string = "Error: Must provide Original & Signature files when enabling Signature mode"
	DeltaFlagsMissingError             string = "Error: Must provide Signature, Updated & Delta files when enabling Delta mode"
	SignatureDeltaFlagsMissingError    string = "Error: Must provide Updated & Delta files when enabling Signature & Delta modes"
	SignatureNotImplementedError       string = "Error: Signature mode not implemented, coming soon"
	DeltaNotImplementedError           string = "Error: Delta mode not implemented, coming soon"
	UnableToCheckFileFolderExistsError string = "Error: Unable to check if file/folder exists"
	FileDoesNotExistError              string = "Error: File does not exist"
	OriginalFileDoesNotExistError      string = "Error: Original file does not exist"
	SearchingForFileButFoundDirError   string = "Error: Searching for a file but found a folder dir"
	OriginalFileIsFolderError          string = "Error: Original file provided is a folder dir"
)
