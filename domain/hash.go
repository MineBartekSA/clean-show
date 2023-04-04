package domain

type Hasher interface {
	Hash(password string) string
	Verify(password, hash string) (bool, error)
}
