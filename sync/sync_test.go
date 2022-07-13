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

func TestGenerateSignature(t *testing.T) {
	t.Run("should return `Signature, nil` when successfully processed all file data for Signature", func(t *testing.T) {
		// Setup
		reader := readerMock{isReadError: false, readSize: int(testChunk)}
		hasReadByte := false
		updatedBuffer := []byte{'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', testBufferNextChar}
		expectedSignature := []models.Signature{{Weak: 76935130210, Strong: "f39dac6cbaba535e2c207cd0cd8f154974223c848f727f98b3564cea569b41cf"}, {Weak: 16426995555, Strong: "2c9d26566889bcb66e96d74b97b14bc36cfd8c2949ab289fff2caeb0422e91b0"}}
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
		require.Equal(t, []models.Signature{}, signature)
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
		require.Equal(t, []models.Signature{}, signature)
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
