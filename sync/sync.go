package sync

import (
	"errors"
	"fmt"
	"io"

	"github.com/curtismenmuir/go-file-diff/constants"
	"github.com/curtismenmuir/go-file-diff/utils"
)

var (
	logger         = utils.Logger
	rollBuffer     = roll
	chunk      int = 10
)

// FileReader interface for mocking bufio.Reader
type Reader interface {
	Read(p []byte) (int, error)
	ReadByte() (byte, error)
}

// GenerateSignature() will create a rolling "hash" buffer and loop through file until it reaches EOF
// NOTE: Signature generation still to come
// Function returns `nil` when successful
// Function returns `error` when unsuccessful
func GenerateSignature(reader Reader) error {
	// Create buffer based on chunk size
	buffer, err := populateBuffer(reader)
	if err != nil {
		return err
	}

	// TODO Generate Week hash
	// TODO Generate Strong hash
	// TODO Store values in Signature Table

	logger(fmt.Sprintf("Initial Buffer = %q", buffer[:]), true)

	// Loop until EOF
	for {
		var initialByte byte = 0
		var nextByte byte = 0
		// Roll buffer to next position
		buffer, initialByte, nextByte, err = rollBuffer(reader, buffer)
		if err != nil {
			// Break loop when EOF returned
			if err.Error() == constants.EndOfFileError {
				break
			}

			// Handle errors
			return err
		}

		// TODO Generate Week hash
		// TODO Generate Strong hash
		// TODO Store values in Signature Table

		logger(fmt.Sprintf("Initial Byte = %c", initialByte), true)
		logger(fmt.Sprintf("New Byte = %c", nextByte), true)
		logger(fmt.Sprintf("Buffer = %q", buffer[:]), true)
	}

	// TODO return Signature table
	return nil
}

// pop() will remove the first item from a provided buffer
// Function returns `updatedBuffer, initialByte`
// Note: initialByte is the item popped from buffer
func pop(buffer []byte) ([]byte, byte) {
	// Create new slice
	buf := make([]byte, 0)
	// Get initial byte
	initialByte := buffer[0]
	// Fill slice with buffer items from position 1 (pop)
	buf = append(buf, buffer[1:]...)
	return buf, initialByte
}

// populateBuffer() will create a new buffer and populate it, based on `chuck` size, from the provided file reader
// Function will return `buffer, nil` when successful
// Function will return `emptyBuffer, EOF` error when reader reaches end of file
// Function will return `emptyBuffer, error` when unable to read from file
func populateBuffer(reader Reader) ([]byte, error) {
	// Create buffer based on chunk size
	buffer := make([]byte, chunk)
	// Fill buffer from file reader
	n, err := reader.Read(buffer)
	if err != nil {
		// Handle EOF error
		if err == io.EOF {
			return []byte{}, errors.New(constants.EndOfFileError)
		}

		return []byte{}, err
	}

	if n == 0 {
		// Handle EOF error
		return []byte{}, errors.New(constants.EndOfFileError)
	}

	return buffer, nil
}

// push() will append the provided byte to the end of the provided buffer
// Function returns `updatedBuffer`
func push(buffer []byte, item byte) []byte {
	return append(buffer, item)
}

// roll() will move the rolling hash function to the next position
// This will include: read next item from file; popping 1st item from buffer; pushing new item to end of buffer;
// Function will return `updatedBuffer, initialByte, nextByte, nil` when successful
// Note: initialByte = byte popped from first position
// Note: nextByte = byte pushed onto end of buffer
// Function will return `emptyBuffer, 0, 0, EOL` error when read EOF from file reader
// Function will return `emptyBuffer, 0, 0, error` when unable to read byte from file
func roll(reader Reader, buffer []byte) ([]byte, byte, byte, error) {
	// Read a byte from file reader
	nextByte, err := reader.ReadByte()
	if err != nil {
		// Handle EOF error
		if err == io.EOF {
			return []byte{}, 0, 0, errors.New(constants.EndOfFileError)
		}

		return []byte{}, 0, 0, err
	}

	// Pop initial byte
	buf, initialByte := pop(buffer)
	// Push new byte
	buf = push(buf, nextByte)
	return buf, initialByte, nextByte, nil
}
