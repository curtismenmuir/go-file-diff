package sync

import (
	"errors"
	"io"
	"testing"

	"github.com/curtismenmuir/go-file-diff/constants"
	"github.com/curtismenmuir/go-file-diff/models"
	"github.com/stretchr/testify/require"
)

var (
	errorMessage          string = "Some Error"
	testChunk             int64  = 16
	testByte              byte   = 5
	testBuffer                   = []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p'}
	testBufferHash               = int64(76935130210)
	testBufferNextChar    byte   = 'q'
	testBufferUpdatedHash        = int64(49921073876)
	testBufferStrongHash  string = "f39dac6cbaba535e2c207cd0cd8f154974223c848f727f98b3564cea569b41cf"
)

// Mock for Reader interface
type readerMock struct {
	// Set test props
	isReadError     bool
	isReadByteError bool
	mockError       error
	readSize        int
}

// Overwrite readerMock.Read() to consider test prop
func (r readerMock) Read(p []byte) (int, error) {
	// Throw error if isReadError set
	if r.isReadError {
		return 0, r.mockError
	}

	return r.readSize, nil
}

// Overwrite readerMock.ReadByte() to consider test prop
func (r readerMock) ReadByte() (byte, error) {
	// Throw error if isReadByteError set
	if r.isReadByteError {
		return 0, r.mockError
	}

	return testByte, nil
}

func TestCompareChecksums(t *testing.T) {
	t.Run("should return `true, item.Head, item.Tail` when weak and strong hashes match block in Signature", func(t *testing.T) {
		// Setup
		expectedHead := 0
		expectedTail := int(testChunk) - 1
		signature := models.Signature{}
		signature[testBufferHash] = models.StrongSignature{Hash: testBufferStrongHash, Head: expectedHead, Tail: expectedTail}
		// Run
		result, head, tail := compareChecksums(signature, testBuffer, testBufferHash, false)
		// Verify
		require.Equal(t, true, result)
		require.Equal(t, expectedHead, head)
		require.Equal(t, expectedTail, tail)
	})

	t.Run("should return `false, -1, -1` when weak hash match block in Signature but not strong hash", func(t *testing.T) {
		// Setup
		buffer := make([]byte, 0)
		buffer = append(buffer, testBuffer[1:]...)
		buffer = append(buffer, '1')
		signature := models.Signature{}
		signature[testBufferHash] = models.StrongSignature{Hash: testBufferStrongHash, Head: 0, Tail: 15}
		// Run
		result, head, tail := compareChecksums(signature, buffer, testBufferHash, false)
		// Verify
		require.Equal(t, false, result)
		require.Equal(t, -1, head)
		require.Equal(t, -1, tail)
	})

	t.Run("should return `false, -1, -1` when weak hash does not match any blocks in Signature", func(t *testing.T) {
		// Setup
		signature := models.Signature{}
		signature[testBufferHash] = models.StrongSignature{Hash: testBufferStrongHash, Head: 0, Tail: 15}
		// Run
		result, head, tail := compareChecksums(signature, testBuffer, 123, false)
		// Verify
		require.Equal(t, false, result)
		require.Equal(t, -1, head)
		require.Equal(t, -1, tail)
	})
}

func TestGenerateDelta(t *testing.T) {
	t.Run("should return `delta, nil` when Updated file contains new block at the beginning of Original file", func(t *testing.T) {
		// Setup
		reader := readerMock{isReadError: false, readSize: int(testChunk)}
		rollCount := 0
		newBlock := []byte{'1', '2', '3'}
		initialBuffer := []byte{newBlock[0], newBlock[1], newBlock[2], 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm'}
		modifiedBlock := []byte{'n', 'o', 'p'}
		// Initialise Signature
		signature := models.Signature{}
		signature[testBufferHash] = models.StrongSignature{Hash: testBufferStrongHash, Head: 0, Tail: 15}
		// Initialise Delta
		expectedDelta := make(models.Delta)
		// Add missing block
		expectedDelta[0] = models.Block{Head: 0, Tail: 2, IsModified: true, Value: newBlock}
		// Add new block
		expectedDelta[3] = models.Block{Head: 0, Tail: 15, IsModified: false, Value: []byte{}}
		// Mock
		initialiseBuffer = func(reader Reader, chunkSize int64) ([]byte, error) {
			return initialBuffer, nil
		}

		rollBuffer = func(reader Reader, buffer []byte) ([]byte, byte, byte, error) {
			// Return EOF to simulate reaching EOF
			if rollCount == len(modifiedBlock) {
				return []byte{}, 0, 0, errors.New(constants.EndOfFileError)
			}

			// Roll buffer
			initialByte := initialBuffer[0]
			nextByte := modifiedBlock[rollCount]
			buf := make([]byte, 0)
			buf = append(buf, initialBuffer[1:]...)
			buf = append(buf, nextByte)
			initialBuffer = buf
			rollCount++
			return initialBuffer, initialByte, nextByte, nil
		}

		// Run
		delta, err := GenerateDelta(reader, signature, false)
		// Verify
		require.Equal(t, nil, err)
		require.Equal(t, len(expectedDelta), len(delta))
		require.Equal(t, expectedDelta, delta)
	})

	t.Run("should return `delta, nil` when Updated file contains new block in middle of Original file", func(t *testing.T) {
		// Setup
		reader := readerMock{isReadError: false, readSize: int(testChunk)}
		rollCount := 0
		initialBuffer := testBuffer
		newBlock := []byte{'q', 'r', 's'}
		modifiedBlock := []byte{'t', 'u', 'v', 'w', 'x', 'y', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j'}
		// Initialise Signature
		signature := models.Signature{}
		signatureBuffer := testBuffer
		head := 0
		tail := int(testChunk) - 1
		signature[generateWeakHash(signatureBuffer, testChunk)] = models.StrongSignature{Hash: generateStrongHash(signatureBuffer, testChunk), Head: head, Tail: tail}
		for index := range modifiedBlock {
			head++
			tail++
			buf := make([]byte, 0)
			buf = append(buf, signatureBuffer[1:]...)
			buf = append(buf, modifiedBlock[index])
			signatureBuffer = buf
			signature[generateWeakHash(signatureBuffer, testChunk)] = models.StrongSignature{Hash: generateStrongHash(signatureBuffer, testChunk), Head: head, Tail: tail}
		}

		// Add new block to modified items for Updated file
		modifiedBlock = append(newBlock[:], modifiedBlock[:]...)
		// Initialise Delta
		expectedDelta := make(models.Delta)
		// Add matched block
		expectedDelta[0] = models.Block{Head: 0, Tail: 15, IsModified: false, Value: []byte{}}
		// Add missing block
		expectedDelta[16] = models.Block{Head: 0, Tail: 2, IsModified: true, Value: newBlock}
		// Add matched block
		expectedDelta[19] = models.Block{Head: 16, Tail: 31, IsModified: false, Value: []byte{}}
		// Mock
		initialiseBuffer = func(reader Reader, chunkSize int64) ([]byte, error) {
			return initialBuffer, nil
		}

		rollBuffer = func(reader Reader, buffer []byte) ([]byte, byte, byte, error) {
			// Return EOF to simulate reaching EOF
			if rollCount == len(modifiedBlock) {
				return []byte{}, 0, 0, errors.New(constants.EndOfFileError)
			}

			// Roll buffer
			initialByte := initialBuffer[0]
			nextByte := modifiedBlock[rollCount]
			buf := make([]byte, 0)
			buf = append(buf, initialBuffer[1:]...)
			buf = append(buf, nextByte)
			initialBuffer = buf
			rollCount++
			return initialBuffer, initialByte, nextByte, nil
		}

		// Run
		delta, err := GenerateDelta(reader, signature, false)
		// Verify
		require.Equal(t, nil, err)
		require.Equal(t, len(expectedDelta), len(delta))
		require.Equal(t, expectedDelta, delta)
	})

	t.Run("should return `delta, nil` when Updated file contains new block at the end of Original file", func(t *testing.T) {
		// Setup
		reader := readerMock{isReadError: false, readSize: int(testChunk)}
		rollCount := 0
		initialBuffer := testBuffer
		modifiedBlock := []byte{'q', 'r', 's'}
		// Initialise Signature
		signature := models.Signature{}
		signature[testBufferHash] = models.StrongSignature{Hash: testBufferStrongHash, Head: 0, Tail: 15}
		// Initialise Delta
		expectedDelta := make(models.Delta)
		// Add matched block
		expectedDelta[0] = models.Block{Head: 0, Tail: 15, IsModified: false, Value: []byte{}}
		// Add missing block
		expectedDelta[16] = models.Block{Head: 0, Tail: 2, IsModified: true, Value: modifiedBlock}
		// Mock
		initialiseBuffer = func(reader Reader, chunkSize int64) ([]byte, error) {
			return initialBuffer, nil
		}

		rollBuffer = func(reader Reader, buffer []byte) ([]byte, byte, byte, error) {
			// Return EOF to simulate reaching EOF
			if rollCount == len(modifiedBlock) {
				return []byte{}, 0, 0, errors.New(constants.EndOfFileError)
			}

			// Roll buffer
			initialByte := initialBuffer[0]
			nextByte := modifiedBlock[rollCount]
			buf := make([]byte, 0)
			buf = append(buf, initialBuffer[1:]...)
			buf = append(buf, nextByte)
			initialBuffer = buf
			rollCount++
			return initialBuffer, initialByte, nextByte, nil
		}

		// Run
		delta, err := GenerateDelta(reader, signature, false)
		// Verify
		require.Equal(t, nil, err)
		require.Equal(t, len(expectedDelta), len(delta))
		require.Equal(t, expectedDelta, delta)
	})

	t.Run("should return `delta, nil` when Updated file contains deleted blocks from Original file", func(t *testing.T) {
		// Setup
		reader := readerMock{isReadError: false, readSize: int(testChunk)}
		rollCount := 0
		initialBuffer := testBuffer
		modifiedBlock := []byte{'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}
		// Initialise Signature
		signature := models.Signature{}
		signatureBuffer := testBuffer
		head := 0
		tail := int(testChunk) - 1
		signature[generateWeakHash(signatureBuffer, testChunk)] = models.StrongSignature{Hash: generateStrongHash(signatureBuffer, testChunk), Head: head, Tail: tail}
		for index := range modifiedBlock {
			head++
			tail++
			buf := make([]byte, 0)
			buf = append(buf, signatureBuffer[1:]...)
			buf = append(buf, modifiedBlock[index])
			signatureBuffer = buf
			signature[generateWeakHash(signatureBuffer, testChunk)] = models.StrongSignature{Hash: generateStrongHash(signatureBuffer, testChunk), Head: head, Tail: tail}
		}

		// Remove first item from modified block to simulate deleted item in Updated file
		modifiedBlock = modifiedBlock[1:]
		// Initialise Delta
		expectedDelta := make(models.Delta)
		// Add matched block
		expectedDelta[0] = models.Block{Head: 0, Tail: 15, IsModified: false, Value: []byte{}}
		// Add matched block
		expectedDelta[16] = models.Block{Head: 17, Tail: 32, IsModified: false, Value: []byte{}}
		// Mock
		initialiseBuffer = func(reader Reader, chunkSize int64) ([]byte, error) {
			return initialBuffer, nil
		}

		rollBuffer = func(reader Reader, buffer []byte) ([]byte, byte, byte, error) {
			if rollCount == len(modifiedBlock) {
				// Return EOF to simulate reaching EOF
				return []byte{}, 0, 0, errors.New(constants.EndOfFileError)
			}

			// Roll buffer
			initialByte := initialBuffer[0]
			nextByte := modifiedBlock[rollCount]
			buf := make([]byte, 0)
			buf = append(buf, initialBuffer[1:]...)
			buf = append(buf, nextByte)
			initialBuffer = buf
			rollCount++
			return initialBuffer, initialByte, nextByte, nil
		}

		// Run
		delta, err := GenerateDelta(reader, signature, false)
		// Verify
		require.Equal(t, nil, err)
		require.Equal(t, len(expectedDelta), len(delta))
		require.Equal(t, expectedDelta, delta)
	})

	t.Run("should return `delta, nil` when Updated file contains numerous updated blocks from Original file", func(t *testing.T) {
		// Setup
		reader := readerMock{isReadError: false, readSize: int(testChunk)}
		rollCount := 0
		var initialModifiedBlock byte = '1'
		var finalModifiedBlock byte = '5'
		initialBuffer := []byte{initialModifiedBlock, 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o'}
		modifiedBlock := []byte{'p', '2', '3', '4', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'a', 'b', 'c', 'd', 'e', 'f', 'g', finalModifiedBlock}
		initialMatchHead := 0
		initialMatchTail := 15
		secondMatchHead := 16
		secondMatchTail := 31
		// Initialise Signature
		signature := models.Signature{}
		signature[testBufferHash] = models.StrongSignature{Hash: testBufferStrongHash, Head: initialMatchHead, Tail: initialMatchTail}
		signature[44661510977] = models.StrongSignature{Hash: "caf9aa676718c8ecebd197c227332bcb342d39187aec40c920f63eab28b6ab87", Head: secondMatchHead, Tail: secondMatchTail}
		// Initialise Delta
		expectedDelta := make(models.Delta)
		// Add missing block
		expectedDelta[0] = models.Block{Head: 0, Tail: 0, IsModified: true, Value: []byte{initialModifiedBlock}}
		// Add matched block
		expectedDelta[1] = models.Block{Head: initialMatchHead, Tail: initialMatchTail, IsModified: false, Value: []byte{}}
		// Add missing block
		expectedDelta[17] = models.Block{Head: 0, Tail: 2, IsModified: true, Value: []byte{modifiedBlock[1], modifiedBlock[2], modifiedBlock[3]}}
		// Add matched block
		expectedDelta[20] = models.Block{Head: secondMatchHead, Tail: secondMatchTail, IsModified: false, Value: []byte{}}
		// Add missing block
		expectedDelta[36] = models.Block{Head: 0, Tail: 0, IsModified: true, Value: []byte{finalModifiedBlock}}
		// Mock
		initialiseBuffer = func(reader Reader, chunkSize int64) ([]byte, error) {
			return initialBuffer, nil
		}

		rollBuffer = func(reader Reader, buffer []byte) ([]byte, byte, byte, error) {
			// Return EOF to simulate reaching EOF
			if rollCount == len(modifiedBlock) {
				return []byte{}, 0, 0, errors.New(constants.EndOfFileError)
			}

			// Roll buffer
			initialByte := initialBuffer[0]
			nextByte := modifiedBlock[rollCount]
			buf := make([]byte, 0)
			buf = append(buf, initialBuffer[1:]...)
			buf = append(buf, nextByte)
			initialBuffer = buf
			rollCount++
			return initialBuffer, initialByte, nextByte, nil
		}

		// Run
		delta, err := GenerateDelta(reader, signature, false)
		// Verify
		require.Equal(t, nil, err)
		require.Equal(t, len(expectedDelta), len(delta))
		require.Equal(t, expectedDelta, delta)
	})

	t.Run("should return `emptyDelta, UpdatedFileHasNoChangesError` when Updated file has no changes from Original", func(t *testing.T) {
		// Setup
		reader := readerMock{isReadError: false, readSize: int(testChunk)}
		hasReadByte := false
		updatedBuffer := []byte{'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', testBufferNextChar}
		signature := models.Signature{}
		signature[testBufferHash] = models.StrongSignature{Hash: testBufferStrongHash, Head: 0, Tail: 15}
		signature[16426995555] = models.StrongSignature{Hash: "2c9d26566889bcb66e96d74b97b14bc36cfd8c2949ab289fff2caeb0422e91b0", Head: 1, Tail: 16}
		expectedError := errors.New(constants.UpdatedFileHasNoChangesError)

		// Mock
		initialiseBuffer = func(reader Reader, chunkSize int64) ([]byte, error) {
			return testBuffer, nil
		}

		rollBuffer = func(reader Reader, buffer []byte) ([]byte, byte, byte, error) {
			if !hasReadByte {
				// Flip hasReadByte to stop ReadByte loop by simulating EOF
				hasReadByte = true
				return updatedBuffer, 1, 5, nil
			}

			return []byte{}, 0, 0, errors.New(constants.EndOfFileError)
		}

		// Run
		delta, err := GenerateDelta(reader, signature, false)
		// Verify
		require.Equal(t, expectedError, err)
		require.Equal(t, models.Delta{}, delta)
	})

	t.Run("should return `emptyDelta, error` when unable to populate buffer from file", func(t *testing.T) {
		// Setup
		reader := readerMock{isReadError: false, readSize: int(testChunk)}
		signature := models.Signature{}
		expectedError := errors.New(errorMessage)
		expectedDelta := models.Delta{}
		// Mock
		initialiseBuffer = func(reader Reader, chunkSize int64) ([]byte, error) {
			return []byte{}, expectedError
		}

		// Run
		delta, err := GenerateDelta(reader, signature, false)
		// Verify
		require.Equal(t, expectedError, err)
		require.Equal(t, expectedDelta, delta)
	})

	t.Run("should return `emptyDelta, error` when unable to read data from Updated file to roll buffer", func(t *testing.T) {
		// Setup
		reader := readerMock{isReadError: false, readSize: int(testChunk)}
		expectedError := errors.New(errorMessage)
		signature := models.Signature{}
		expectedDelta := models.Delta{}
		// Mock
		initialiseBuffer = func(reader Reader, chunkSize int64) ([]byte, error) {
			return testBuffer, nil
		}

		rollBuffer = func(reader Reader, buffer []byte) ([]byte, byte, byte, error) {
			return []byte{}, 0, 0, expectedError
		}

		// Run
		delta, err := GenerateDelta(reader, signature, false)
		// Verify
		require.Equal(t, expectedError, err)
		require.Equal(t, expectedDelta, delta)
	})
}

func TestGenerateMatchedBlock(t *testing.T) {
	t.Run("should return `matchingBlock, blockHead, initialBlockMatches` after increasing block tail position when already processing a matching block (EG previous roll matched)", func(t *testing.T) {
		// Setup
		delta := models.Delta{}
		exists := true
		initialBlockMatches := true
		blockHead := 0
		deltaHead := 1
		rollHead := 2
		rollTail := 17
		rollExists := initialBlockMatches
		block := models.Block{Head: blockHead, Tail: blockHead, IsModified: false, Value: []byte{}}
		expectedBlock := models.Block{Head: blockHead, Tail: blockHead + 1, IsModified: false, Value: []byte{}}
		expectedInitialBlockMatches := true
		expectedBlockHead := 0
		// Run
		block, blockHead, initialBlockMatches = generateMatchedBlock(delta, block, exists, initialBlockMatches, blockHead, deltaHead, rollHead, rollTail, rollExists, false)
		// Verify
		require.Equal(t, 0, len(delta))
		require.Equal(t, expectedBlock, block)
		require.Equal(t, expectedBlockHead, blockHead)
		require.Equal(t, expectedInitialBlockMatches, initialBlockMatches)
	})

	t.Run("should return `matchingBlock, blockHead, initialBlockMatches` after adding initial missing block to Delta (EG new block added to start of file)", func(t *testing.T) {
		// Setup
		delta := models.Delta{}
		exists := false
		initialBlockMatches := false
		blockHead := 0
		deltaHead := 1
		rollHead := 2
		rollTail := 17
		rollExists := true
		value := []byte{'a'}
		block := models.Block{Head: blockHead, Tail: blockHead, IsModified: true, Value: value}
		expectedBlock := models.Block{Head: rollHead, Tail: rollTail, IsModified: false, Value: []byte{}}
		expectedInitialBlockMatches := !initialBlockMatches
		expectedBlockHead := 1
		// Run
		block, blockHead, initialBlockMatches = generateMatchedBlock(delta, block, exists, initialBlockMatches, blockHead, deltaHead, rollHead, rollTail, rollExists, false)
		// Verify
		require.Equal(t, 1, len(delta))
		require.Equal(t, value, delta[0].Value)
		require.Equal(t, expectedBlock, block)
		require.Equal(t, expectedBlockHead, blockHead)
		require.Equal(t, expectedInitialBlockMatches, initialBlockMatches)
	})

	t.Run("should return `matchingBlock, blockHead, initialBlockMatches` after reducing block to contain only missing chars and adding to Delta (EG new block after initial matched block)", func(t *testing.T) {
		// Setup
		delta := models.Delta{}
		exists := false
		initialBlockMatches := true
		blockHead := 0
		blockTail := 16
		deltaHead := 1
		rollHead := 2
		rollTail := 17
		rollExists := true
		value := testBuffer
		value = append(value, testBufferNextChar)
		block := models.Block{Head: blockHead, Tail: blockTail, IsModified: true, Value: value}
		expectedValue := []byte{testBuffer[0], testBuffer[1]}
		expectedBlock := models.Block{Head: rollHead, Tail: rollTail, IsModified: false, Value: []byte{}}
		expectedInitialBlockMatches := initialBlockMatches
		expectedBlockHead := 1
		// Run
		block, blockHead, initialBlockMatches = generateMatchedBlock(delta, block, exists, initialBlockMatches, blockHead, deltaHead, rollHead, rollTail, rollExists, false)
		// Verify
		require.Equal(t, 1, len(delta))
		require.Equal(t, expectedValue, delta[0].Value)
		require.Equal(t, expectedBlock, block)
		require.Equal(t, expectedBlockHead, blockHead)
		require.Equal(t, expectedInitialBlockMatches, initialBlockMatches)
	})
}

func TestGenerateMissingBlock(t *testing.T) {
	t.Run("should return `missingBlock, blockHead` after adding previous matched block to Delta (EG found missing block after previous roll matched)", func(t *testing.T) {
		// Setup
		delta := models.Delta{}
		exists := true
		initialBlockMatches := true
		blockHead := 0
		nextByte := testBufferNextChar
		buffer := testBuffer
		expectedValue := []byte{}
		block := models.Block{Head: blockHead, Tail: blockHead, IsModified: false, Value: expectedValue}
		expectedBlock := models.Block{Head: blockHead, Tail: blockHead, IsModified: exists, Value: []byte{nextByte}}
		expectedBlockHead := 1
		// Run
		block, blockHead = generateMissingBlock(delta, block, exists, initialBlockMatches, blockHead, nextByte, buffer, false)
		// Verify
		require.Equal(t, 1, len(delta))
		require.Equal(t, expectedValue, delta[0].Value)
		require.Equal(t, expectedBlock, block)
		require.Equal(t, expectedBlockHead, blockHead)
	})

	t.Run("should return `missingBlock, blockHead` after adding first buffer to missing block and incrementing tail position when processing initial missing block (EG new block added to start of file)", func(t *testing.T) {
		// Setup
		delta := models.Delta{}
		exists := false
		initialBlockMatches := false
		blockHead := 0
		nextByte := testBufferNextChar
		buffer := testBuffer
		block := models.Block{Head: blockHead, Tail: blockHead, IsModified: true, Value: []byte{testBufferNextChar}}
		expectedBlock := models.Block{Head: blockHead, Tail: blockHead + 1, IsModified: true, Value: []byte{testBufferNextChar, buffer[0]}}
		expectedBlockHead := 0
		// Run
		block, blockHead = generateMissingBlock(delta, block, exists, initialBlockMatches, blockHead, nextByte, buffer, false)
		// Verify
		require.Equal(t, 0, len(delta))
		require.Equal(t, expectedBlock, block)
		require.Equal(t, expectedBlockHead, blockHead)
	})

	t.Run("should return `missingBlock, blockHead` after adding new byte to missing block and incrementing tail position when processing missing block after initial matching block (EG new block after initial matched block)", func(t *testing.T) {
		// Setup
		delta := models.Delta{}
		exists := false
		initialBlockMatches := true
		blockHead := 0
		buffer := testBuffer
		nextByte := testBuffer[15]
		block := models.Block{Head: blockHead, Tail: blockHead, IsModified: true, Value: []byte{testBufferNextChar}}
		expectedBlock := models.Block{Head: blockHead, Tail: blockHead + 1, IsModified: true, Value: []byte{testBufferNextChar, nextByte}}
		expectedBlockHead := 0
		// Run
		block, blockHead = generateMissingBlock(delta, block, exists, initialBlockMatches, blockHead, nextByte, buffer, false)
		// Verify
		require.Equal(t, 0, len(delta))
		require.Equal(t, expectedBlock, block)
		require.Equal(t, expectedBlockHead, blockHead)
	})
}

func TestGenerateSignature(t *testing.T) {
	t.Run("should return `Signature, nil` when successfully processed all file data for Signature", func(t *testing.T) {
		// Setup
		reader := readerMock{isReadError: false, readSize: int(testChunk)}
		hasReadByte := false
		updatedBuffer := []byte{'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', testBufferNextChar}
		expectedSignature := models.Signature{}
		expectedSignature[testBufferHash] = models.StrongSignature{Hash: testBufferStrongHash, Head: 0, Tail: 15}
		expectedSignature[16426995555] = models.StrongSignature{Hash: "2c9d26566889bcb66e96d74b97b14bc36cfd8c2949ab289fff2caeb0422e91b0", Head: 1, Tail: 16}
		// Mock
		initialiseBuffer = func(reader Reader, chunkSize int64) ([]byte, error) {
			return testBuffer, nil
		}

		rollBuffer = func(reader Reader, buffer []byte) ([]byte, byte, byte, error) {
			if !hasReadByte {
				// Flip hasReadByte to stop ReadByte loop by simulating EOF
				hasReadByte = true
				return updatedBuffer, 1, 5, nil
			}

			return []byte{}, 0, 0, errors.New(constants.EndOfFileError)
		}

		// Run
		signature, err := GenerateSignature(reader, false)
		// Verify
		require.Equal(t, nil, err)
		require.NotEqual(t, nil, signature)
		require.Equal(t, expectedSignature, signature)
	})

	t.Run("should return `emptySignature, error` when unable to populate buffer from file", func(t *testing.T) {
		// Setup
		expectedError := errors.New(errorMessage)
		reader := readerMock{isReadError: true, mockError: expectedError}
		// Mock
		initialiseBuffer = func(reader Reader, chunkSize int64) ([]byte, error) {
			return []byte{}, expectedError
		}

		// Run
		signature, err := GenerateSignature(reader, false)
		// Verify
		require.Equal(t, expectedError, err)
		require.Equal(t, models.Signature{}, signature)
	})

	t.Run("should return `emptySignature, error` when unable to read data from file to roll buffer", func(t *testing.T) {
		// Setup
		expectedError := errors.New(errorMessage)
		reader := readerMock{isReadError: false, readSize: int(testChunk)}
		// Mock
		initialiseBuffer = func(reader Reader, chunkSize int64) ([]byte, error) {
			return testBuffer, nil
		}

		rollBuffer = func(reader Reader, buffer []byte) ([]byte, byte, byte, error) {
			return []byte{}, 0, 0, expectedError
		}

		// Run
		signature, err := GenerateSignature(reader, false)
		// Verify
		require.Equal(t, expectedError, err)
		require.Equal(t, models.Signature{}, signature)
	})
}

func TestGenerateWeakHash(t *testing.T) {
	t.Run("should return a consistent `resultHash` after hashing the provided buffer", func(t *testing.T) {
		// Run
		resultHash := generateWeakHash(testBuffer, testChunk)
		// Verify
		require.Equal(t, testBufferHash, resultHash)
	})

	t.Run("should generate a different `resultHash` for different buffers", func(t *testing.T) {
		// Setup
		buffer := []byte{'f', 'b', 'c', 'e', 'e', 'f', 'g', 4, 'i', 2, 'k', 'l', 'm', '£', 'o', 'p'}
		// Run
		resultHash := generateWeakHash(testBuffer, testChunk)
		differentHash := generateWeakHash(buffer, testChunk)
		// Verify
		require.Equal(t, testBufferHash, resultHash)
		require.NotEqual(t, differentHash, resultHash)
	})

	t.Run("should generate a different `resultHash` for hashes which have been reversed (EG byte order important)", func(t *testing.T) {
		// Setup
		buffer := []byte{'p', 'o', 'n', 'm', 'l', 'k', 'j', 'i', 'h', 'g', 'f', 'e', 'd', 'c', 'b', 'a'}
		anotherBuffer := []byte{'b', 'a', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p'}
		// Run
		resultHash := generateWeakHash(testBuffer, testChunk)
		differentHash := generateWeakHash(buffer, testChunk)
		anotherHash := generateWeakHash(anotherBuffer, testChunk)
		// Verify
		require.Equal(t, testBufferHash, resultHash)
		require.NotEqual(t, differentHash, resultHash)
		require.NotEqual(t, differentHash, anotherHash)
		require.NotEqual(t, anotherHash, resultHash)
	})

	t.Run("should generate a valid `resultHash` when using max byte size (255)", func(t *testing.T) {
		// Setup
		buffer := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
		expectedHash := int64(11415635451)
		// Run
		resultHash := generateWeakHash(buffer, testChunk)
		// Verify
		require.Equal(t, expectedHash, resultHash)
	})
}

func TestGenerateStrongHash(t *testing.T) {
	t.Run("should return SHA-256 `hash` of provided buffer as a Hex string", func(t *testing.T) {
		// Run
		hash := generateStrongHash(testBuffer, testChunk)
		// Verify
		require.Equal(t, testBufferStrongHash, hash)
	})
}

func TestModulo(t *testing.T) {
	t.Run("should return `result` when calculated the remainder between 2 values (mod)", func(t *testing.T) {
		// Setup
		x := int64(10)
		y := int64(4)
		expectedResult := int64(2)
		// Run
		result := modulo(x, y)
		// Verify
		require.Equal(t, expectedResult, result)
	})

	t.Run("should implement Euclidean modulus, which differs from Go's mod operator", func(t *testing.T) {
		// Setup
		x := int64(-10)
		y := int64(4)
		expectedResult := int64(2)
		goModResult := x % y
		// Run
		result := modulo(x, y)
		// Verify
		require.Equal(t, expectedResult, result)
		require.NotEqual(t, goModResult, result)
	})
}

func TestPop(t *testing.T) {
	t.Run("should return `new byte[], initialByte` with initial byte popped from beginning of array", func(t *testing.T) {
		// Setup
		var expectedByte byte = 0
		initialBuffer := []byte{expectedByte, 1, 2, 3, 4}
		expectedBuffer := []byte{1, 2, 3, 4}
		// Run
		buffer, initialByte := pop(initialBuffer)
		// Verify
		require.Equal(t, expectedBuffer, buffer)
		require.Equal(t, expectedByte, initialByte)
	})
}

func TestPopulateBuffer(t *testing.T) {
	t.Run("should return `buffer, nil` when successfully populated buffer from file", func(t *testing.T) {
		// Setup
		reader := readerMock{isReadError: false, readSize: int(testChunk)}
		expectedBuffer := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		// Run
		buffer, err := populateBuffer(reader, testChunk)
		// Verify
		require.Equal(t, expectedBuffer, buffer)
		require.Equal(t, nil, err)
	})

	t.Run("should return `emptyBuffer, EOF error` reader returns EOF error", func(t *testing.T) {
		// Setup
		expectedError := errors.New(constants.EndOfFileError)
		reader := readerMock{isReadError: false, readSize: 0}
		expectedBuffer := []byte{}
		// Run
		buffer, err := populateBuffer(reader, testChunk)
		// Verify
		require.Equal(t, expectedBuffer, buffer)
		require.Equal(t, expectedError, err)
	})

	t.Run("should return `emptyBuffer, EOF error` reader returns 0 bytes read", func(t *testing.T) {
		// Setup
		expectedError := errors.New(constants.EndOfFileError)
		reader := readerMock{isReadError: true, mockError: io.EOF}
		expectedBuffer := []byte{}
		// Run
		buffer, err := populateBuffer(reader, testChunk)
		// Verify
		require.Equal(t, expectedBuffer, buffer)
		require.Equal(t, expectedError, err)
	})

	t.Run("should return `emptyBuffer, error` reader fails to read data from file", func(t *testing.T) {
		// Setup
		expectedError := errors.New(errorMessage)
		reader := readerMock{isReadError: true, mockError: expectedError}
		expectedBuffer := []byte{}
		// Run
		buffer, err := populateBuffer(reader, testChunk)
		// Verify
		require.Equal(t, expectedBuffer, buffer)
		require.Equal(t, expectedError, err)
	})
}

func TestPush(t *testing.T) {
	t.Run("should return new byte[] with new item appended", func(t *testing.T) {
		// Setup
		initialBuffer := []byte{0, 1, 2, 3, 4}
		var item byte = 5
		expectedResult := append(initialBuffer, item)
		// Run
		result := push(initialBuffer, item)
		// Verify
		require.Equal(t, expectedResult, result)
	})
}

func TestRoll(t *testing.T) {
	t.Run("should return `updatedBuffer, initialByte, nextByte, nil` when successfully rolled buffer to next position", func(t *testing.T) {
		// Setup
		reader := readerMock{isReadByteError: false}
		var expectedInitByte byte = 0
		initialBuffer := []byte{expectedInitByte, 1, 2, 3, 4}
		expectedBuffer := []byte{1, 2, 3, 4, testByte}
		// Run
		buffer, initialByte, nextByte, err := roll(reader, initialBuffer)
		// Verify
		require.Equal(t, expectedBuffer, buffer)
		require.Equal(t, expectedInitByte, initialByte)
		require.Equal(t, testByte, nextByte)
		require.Equal(t, nil, err)
	})

	t.Run("should return `emptyBuffer, emptyByte, emptyByte, EOF error` when unable to roll to next position as reached EOF", func(t *testing.T) {
		// Setup
		expectedError := errors.New(constants.EndOfFileError)
		reader := readerMock{isReadByteError: true, mockError: io.EOF}
		initialBuffer := []byte{0, 1, 2, 3, 4}
		var expectedByte byte = 0
		// Run
		buffer, initialByte, nextByte, err := roll(reader, initialBuffer)
		// Verify
		require.Equal(t, []byte{}, buffer)
		require.Equal(t, expectedByte, initialByte)
		require.Equal(t, expectedByte, nextByte)
		require.Equal(t, expectedError, err)
	})

	t.Run("should return `emptyBuffer, emptyByte, emptyByte, error` when unable to roll to next position", func(t *testing.T) {
		// Setup
		expectedError := errors.New(errorMessage)
		reader := readerMock{isReadByteError: true, mockError: expectedError}
		initialBuffer := []byte{0, 1, 2, 3, 4}
		var expectedByte byte = 0
		// Run
		buffer, initialByte, nextByte, err := roll(reader, initialBuffer)
		// Verify
		require.Equal(t, []byte{}, buffer)
		require.Equal(t, expectedByte, initialByte)
		require.Equal(t, expectedByte, nextByte)
		require.Equal(t, expectedError, err)
	})
}

func TestRollWeakHash(t *testing.T) {
	t.Run("should return a consistent `updatedHash` after rolling hash to next position", func(t *testing.T) {
		// Run
		result := rollWeakHash(testBufferHash, testBuffer[0], testBufferNextChar, testChunk)
		// Verify
		require.NotEqual(t, testBufferHash, result)
		require.Equal(t, testBufferUpdatedHash, result)
	})

	t.Run("should return an `updatedHash` which matches generating hash with full buffer (eg rolls to correct hash)", func(t *testing.T) {
		// Setup
		buffer := []byte{'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', testBufferNextChar}
		// Run
		result := rollWeakHash(testBufferHash, testBuffer[0], testBufferNextChar, testChunk)
		expectedResult := generateWeakHash(buffer, testChunk)
		// Verify
		require.NotEqual(t, testBufferHash, result)
		require.Equal(t, testBufferUpdatedHash, result)
		require.Equal(t, expectedResult, result)
	})
}
