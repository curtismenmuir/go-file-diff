package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Run("should call log function when verbose flag set to true", func(t *testing.T) {
		// Setup
		invoked := false
		// Mock
		log = func(a ...interface{}) (n int, err error) {
			invoked = true
			return 0, nil
		}
		// Run
		Logger("Some Message", true)
		// Verify results
		require.Equal(t, true, invoked)
	})

	t.Run("should not call log function when verbose flag set to false", func(t *testing.T) {
		// Setup
		invoked := false
		// Mock
		log = func(a ...interface{}) (n int, err error) {
			invoked = true
			return 0, nil
		}
		// Run
		Logger("Some Message", false)
		// Verify
		require.Equal(t, false, invoked)
	})
}
