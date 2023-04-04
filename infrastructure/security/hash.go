package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/minebarteksa/clean-show/domain"
	. "github.com/minebarteksa/clean-show/logger"
	"golang.org/x/crypto/argon2"
)

const (
	// NOTE: These values should be adjusted per deployment since they indicate how fast the hashing process will take
	// The following values were configured for a premium home laptop. Your execution time may vary
	TIME        = uint32(16)
	MEMORY      = uint32(256 * 1024)
	THREADS     = uint8(8)
	KEY_LENGTH  = uint32(512)
	SALT_LENGTH = uint32(128)
)

type argon2idHasher struct {
	version string
}

func NewArgon2idHasher() domain.Hasher {
	return &argon2idHasher{
		version: fmt.Sprintf("v=%d", argon2.Version),
	}
}

func (ah *argon2idHasher) Hash(password string) string {
	salt := generateSalt()
	hash := argon2.IDKey([]byte(password), salt, TIME, MEMORY, THREADS, KEY_LENGTH)
	Log.Debugln("Hash len:", len(hash))

	encSalt := base64.StdEncoding.EncodeToString(salt)
	Log.Debugln("New Salt enc len:", len(encSalt))
	encHash := base64.StdEncoding.EncodeToString(hash)
	Log.Debugln("Hash enc len:", len(encHash))

	hashString := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, MEMORY, TIME, THREADS, encSalt, encHash)
	Log.Debugln("Hash string len:", len(hashString))
	return hashString
}

func (ah *argon2idHasher) Verify(password, hash string) (bool, error) {
	split := strings.SplitN(hash, "$", 6)
	if len(split) < 6 {
		return false, fmt.Errorf("invalid hash")
	}
	if split[2] != ah.version {
		return false, fmt.Errorf("incompatible version")
	}

	time := uint32(0)
	memory := uint32(0)
	threads := uint8(0)
	n, err := fmt.Sscanf(split[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return false, fmt.Errorf("scanf fail err %s", err)
	} else if n != 3 {
		return false, fmt.Errorf("read invalid hash")
	}

	salt, err := base64.StdEncoding.DecodeString(split[4])
	if err != nil {
		return false, err
	}
	rawHash, err := base64.StdEncoding.DecodeString(split[5])
	if err != nil {
		return false, err
	}

	testHash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(rawHash)))
	return subtle.ConstantTimeCompare(rawHash, testHash) == 1, nil
}

func generateSalt() []byte {
	buffer := make([]byte, SALT_LENGTH)
	_, err := rand.Read(buffer)
	if err != nil {
		Log.Panicw("failed to generate new salt", "err", err)
	}
	return buffer
}
