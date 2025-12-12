package server

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewServerPortValidation(t *testing.T) {
	sc := []struct {
		name          string
		restPort      string
		gRPCPort      string
		expectedError error
	}{
		{
			name:          "emprty restPort",
			expectedError: errEmptyPort,
		},
		{
			name:          "string restPort",
			restPort:      "INVALID",
			expectedError: errPortNotAValidNumber,
		},
		{
			name:          "reserved restPort",
			restPort:      "1023",
			expectedError: errPortReserved,
		},
		{
			name:          "beyond range",
			restPort:      "65536",
			expectedError: errPortBeyondRange,
		},
		{
			name:     "valid restPort",
			restPort: "8080",
			gRPCPort: "8081",
		},
	}
	logger := logrus.New()
	for _, tc := range sc {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewServer(logger, tc.restPort, tc.gRPCPort)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

// mock http.Server to replace ListenAndServe and Shutdown
// since http.Server is a struct, we'll simulate the error via goroutine and channels
func TestServer_Run_Error(t *testing.T) {
	s, err := NewServer(logrus.New(), "8099", "8089")
	assert.NoError(t, err)
	//var buf bytes.Buffer
	// Simulate error from ListenAndServe
	go func() {
		time.Sleep(10 * time.Millisecond)
		s.ServerError <- errors.New("listen error")
	}()
	err = s.Run()
	assert.EqualError(t, err, "listen error")
}

func TestServer_Run_Shutdown(t *testing.T) {
	s, err := NewServer(logrus.New(), "8090", "8091")
	assert.NoError(t, err)
	// Simulate shutdown signal after a short delay
	go func() {
		time.Sleep(10 * time.Millisecond)
		s.ShutDown <- os.Kill
	}()
	err = s.Run()
	assert.NoError(t, err)
}
