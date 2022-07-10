package sync

import (
	"errors"
	"io"
	"testing"

	"github.com/curtismenmuir/go-file-diff/constants"
	"github.com/stretchr/testify/require"
)

var (
	errorMessage string = "Some Error"
	testChunk    int    = 10
	testByte     byte   = 5
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
	t.Run("should return `nil` when successfully processed all file data for Signature", func(t *testing.T) {
		// Setup
		reader := readerMock{isReadError: false, readSize: testChunk}
		hasReadByte := false
		// Mock
		rollBuffer = func(reader Reader, buffer []byte) ([]byte, byte, byte, error) {
			if !hasReadByte {
				// Flip hasReadByte to stop ReadByte loop by simulating EOF
				hasReadByte = true
				return []byte{2, 3, 4, 5}, 1, 5, nil
			}

			return []byte{}, 0, 0, errors.New(constants.EndOfFileError)
		}

		// Run
		result := GenerateSignature(reader)
		// Verify
		require.Equal(t, nil, result)
	})

	t.Run("should return `error` when unable to populate buffer from file", func(t *testing.T) {
		// Setup
		expectedError := errors.New(errorMessage)
		reader := readerMock{isReadError: true, mockError: expectedError}
		// Run
		result := GenerateSignature(reader)
		// Verify
		require.Equal(t, expectedError, result)
	})

	t.Run("should return `error` when unable to read data from file to roll buffer", func(t *testing.T) {
		// Setup
		expectedError := errors.New(errorMessage)
		reader := readerMock{isReadError: false, readSize: testChunk}
		// Mock
		rollBuffer = func(reader Reader, buffer []byte) ([]byte, byte, byte, error) {
			return []byte{}, 0, 0, expectedError
		}

		// Run
		result := GenerateSignature(reader)
		// Verify
		require.Equal(t, expectedError, result)
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
		reader := readerMock{isReadError: false, readSize: testChunk}
		expectedBuffer := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		// Run
		buffer, err := populateBuffer(reader)
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
		buffer, err := populateBuffer(reader)
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
		buffer, err := populateBuffer(reader)
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
		buffer, err := populateBuffer(reader)
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
