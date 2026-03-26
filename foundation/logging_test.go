package foundation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger()
	require.NotNil(t, logger)
}

func TestValidatePort(t *testing.T) {
	scenarios := []struct {
		name        string
		port        string
		expectedErr error
	}{
		{name: "empty port", port: "", expectedErr: errEmptyPort},
		{name: "not a number", port: "abc", expectedErr: errPortNotAValidNumber},
		{name: "reserved port", port: "80", expectedErr: errPortReserved},
		{name: "port beyond range", port: "99999", expectedErr: errPortBeyondRange},
		{name: "valid port", port: "8080", expectedErr: nil},
		{name: "lower boundary valid", port: "1024", expectedErr: nil},
		{name: "upper boundary valid", port: "65535", expectedErr: nil},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			err := ValidatePort(sc.port)
			assert.Equal(t, sc.expectedErr, err)
		})
	}
}
