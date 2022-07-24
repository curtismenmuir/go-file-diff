package models

type CMD struct {
	Verbose       bool   `json:"verbose"`
	SignatureMode bool   `json:"signatureMode"`
	DeltaMode     bool   `json:"deltaMode"`
	OriginalFile  string `json:"originalFile"`
	SignatureFile string `json:"signatureFile"`
	UpdatedFile   string `json:"updatedFile"`
	DeltaFile     string `json:"deltaFile"`
}

type StrongSignature struct {
	Hash string `json:"hash"`
	Head int    `json:"head"`
	Tail int    `json:"tail"`
}

type Signature map[int64]StrongSignature

type Block struct {
	Head       int    `json:"head"`
	Tail       int    `json:"tail"`
	IsModified bool   `json:"isModified"`
	Value      []byte `json:"value"`
}

type Delta map[int]Block
