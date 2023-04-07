package security_test

import (
	"crypto/rand"
	"encoding/base64"
	"testing"
	"time"

	"github.com/minebarteksa/clean-show/infrastructure/security"
	"github.com/minebarteksa/clean-show/logger"
	"github.com/stretchr/testify/assert"
)

func TestArgon2idHasher(t *testing.T) { // TODO: Adjust params
	logger.InitProduction()
	hasher := security.NewArgon2idHasher()

	for i := 0; i < 5; i++ {
		password, err := generatePassword()
		assert.NoError(t, err)

		start := time.Now()
		hash := hasher.Hash(password)
		took := time.Since(start)
		assert.NotEmpty(t, hash)
		assert.Less(t, took, time.Second+100*time.Millisecond)
		assert.Greater(t, took, 450*time.Millisecond)
		logger.Log.Infoln("Hash took:", took)

		start = time.Now()
		verified, err := hasher.Verify(password, hash)
		took = time.Since(start)
		assert.NoError(t, err)
		assert.True(t, verified)
		assert.Less(t, took, time.Second+100*time.Millisecond)
		assert.Greater(t, took, 450*time.Millisecond)
		logger.Log.Infoln("Verify took:", took)
	}
}

func generatePassword() (string, error) {
	bytes := make([]byte, 5)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}
