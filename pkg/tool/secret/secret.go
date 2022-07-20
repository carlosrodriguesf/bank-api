//go:generate mockgen -source=${GOFILE} -package=${GOPACKAGE} -destination=${GOPACKAGE}_mock.go

package secret

import (
	"crypto/sha512"
	"fmt"
	"math/rand"
	"time"
)

type (
	Secret interface {
		GenSalt() (hash string)
		Encode(password, salt string) (hash string)
		Verify(decoded, encoded, salt string) (isValid bool)
	}
	secretImpl struct {
	}
)

func New() Secret {
	return &secretImpl{}
}

func (s *secretImpl) Encode(password, salt string) string {
	return hash(password + salt)
}

func (s *secretImpl) Verify(decoded, encoded, salt string) bool {
	return encoded == s.Encode(decoded, salt)
}

func (s *secretImpl) GenSalt() string {
	b := make([]byte, 16)
	rand.Seed(time.Now().UnixNano())
	rand.Read(b)
	r := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return hash(r)
}

func hash(s string) string {
	h := sha512.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
