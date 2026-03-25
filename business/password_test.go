package business

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestGeneratePassword(t *testing.T) {
	scenarios := []struct {
		name      string
		minLength int
		unique    bool
	}{
		{name: "meets minimum length", minLength: 15},
		{name: "produces unique values", unique: true},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			p1, err := GeneratePassword()
			require.NoError(t, err)

			if sc.minLength > 0 {
				assert.GreaterOrEqual(t, len(p1), sc.minLength)
			}

			if sc.unique {
				p2, err := GeneratePassword()
				require.NoError(t, err)
				assert.NotEqual(t, p1, p2)
			}
		})
	}
}

func TestEncryptPassword(t *testing.T) {
	scenarios := []struct {
		name            string
		password        string
		expectValidHash bool
		expectUnique    bool
	}{
		{
			name:            "hashes plain password",
			password:        "myS3cretPass",
			expectValidHash: true,
		},
		{
			name:         "same input produces different hashes due to bcrypt salt",
			password:     "myS3cretPass",
			expectUnique: true,
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			hashed, err := EncryptPassword(sc.password)
			require.NoError(t, err)
			assert.NotEmpty(t, hashed)
			assert.NotEqual(t, sc.password, hashed)

			if sc.expectValidHash {
				err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(sc.password))
				assert.NoError(t, err)
			}

			if sc.expectUnique {
				hashed2, err := EncryptPassword(sc.password)
				require.NoError(t, err)
				assert.NotEqual(t, hashed, hashed2)
			}
		})
	}
}