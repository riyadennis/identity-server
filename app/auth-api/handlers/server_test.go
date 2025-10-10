package handlers

import (
	"log"
	"os"
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
			// using select to listen to channel
			select {
			case err := <-se.ServerError:
				assert.Equal(t, tc.expectedError, err)
			default:
				t.Log("reached default for test " + tc.name)
				close(se.ServerError)
			}
		})
	}
}

func TestServerRun(t *testing.T) {
	s := NewServer("8080")
	logger := log.New(os.Stdout, "IDENTITY-TEST", log.LstdFlags)
	s.Run(nil, nil, logger)
	close(s.ServerError)
	close(s.ShutDown)
}

func TestServerShutDown(t *testing.T) {
	s := NewServer("8080")
	logger := log.New(os.Stdout, "IDENTITY-TEST", log.LstdFlags)
	s.Run(nil, nil, logger)
	s.ShutDown <- os.Interrupt
	close(s.ServerError)
	close(s.ShutDown)
}
