package models

// CMD type.
// This will contain the CMD Flags set by user.
type CMD struct {
	Verbose       bool   `json:"verbose"`
	SignatureMode bool   `json:"signatureMode"`
	DeltaMode     bool   `json:"deltaMode"`
	OriginalFile  string `json:"originalFile"`
	SignatureFile string `json:"signatureFile"`
	UpdatedFile   string `json:"updatedFile"`
	DeltaFile     string `json:"deltaFile"`
}

// StrongSignature type.
// This will be used to contain a SHA-256 hash of the block of data, as well as the Head and Tail position of the bytes in the Original file (EG position of first + last characters).
// EG: StrongSignature{Hash: "some-strong-hash", Head: 0, Tail: 15}.
type StrongSignature struct {
	Hash string `json:"hash"`
	Head int    `json:"head"`
	Tail int    `json:"tail"`
}

// Signature type.
// Items will be indexed by their weak hash.
// EG:
// signature[123]{Hash: "some-strong-hash", Head: 0, Tail: 15}.
// signature[456]{Hash: "another-strong-hash", Head: 0, Tail: 15}.
type Signature map[int64]StrongSignature

// Block type.
// This will be used to store the data for each block to be written to final output file (after patch).
// A matching block from Signature file will use Head + Tail to define the blocks position within the Signature file (EG position of first + last characters).
// EG: Block{Head: 0, Tail: 4, IsModified: false, Value: []bytes{}}.
// A missing block from Signature file will use Value to define the byte array to be added to recreate the Updated file.
// EG: Block{Head: 0, Tail: 4, IsModified: true, Value: []bytes{'a', 'b', 'c', 'd', 'e'}}.

type Block struct {
	Head       int    `json:"head"`
	Tail       int    `json:"tail"`
	IsModified bool   `json:"isModified"`
	Value      []byte `json:"value"`
}

// Delta type.
// Items will be indexed by their position in the final output file.
// EG:
// delta[0]{Head: 0, Tail: 4, IsModified: true, Value: []bytes{'a', 'b', 'c', 'd', 'e'}}.
// delta[5]{Head: 0, Tail: 4, IsModified: false, Value: []bytes{}}.
type Delta map[int]Block
