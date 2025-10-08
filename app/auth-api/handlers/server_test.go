package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServerPortValdation(t *testing.T) {
	sc := []struct {
		name          string
		port          string
		expectedError error
	}{
		{
			name:          "emprty port",
			expectedError: errEmptyPort,
		},
		{
			name:          "string port",
			port:          "INVALID",
			expectedError: errPortNotAValidNumber,
		},
		{
			name:          "reserved port",
			port:          "1023",
			expectedError: errPortReserved,
		},
		{
			name:          "beyond range",
			port:          "65536",
			expectedError: errPortBeyondRange,
		},
		{
			name: "valid port",
			port: "8080",
		},
	}
	for _, tc := range sc {
		t.Run(tc.name, func(t *testing.T) {
			se := NewServer(tc.port)
			err := <-se.serverError
			assert.Equal(t, tc.expectedError, err)
		})
	}

}
