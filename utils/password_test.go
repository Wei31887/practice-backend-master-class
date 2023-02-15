package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashedPassword(t *testing.T) {
	passwd := RandomString(6)
	hashedPd, err := HashedPassword(passwd)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPd)

	// validate the password and hased password
	err = ComparePassword(hashedPd, passwd)
	require.NoError(t, err)

	passwd2 := RandomString(6)
	err = ComparePassword(passwd2, passwd)
	require.Error(t, err)
}