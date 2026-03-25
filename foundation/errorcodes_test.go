package foundation

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomError_Error(t *testing.T) {
	err := &CustomError{
		Code: ValidationFailed,
		Err:  errors.New("field required"),
	}
	assert.Equal(t, "field required", err.Error())
}