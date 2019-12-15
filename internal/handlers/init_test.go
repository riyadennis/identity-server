package handlers

import "testing"

func TestGeneratePassword(t *testing.T) {
	pass, _ := generatePassword()
	t.Log(pass)
}
