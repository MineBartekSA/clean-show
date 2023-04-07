package test

import "fmt"

type MockHasher struct {
	HashMap map[string]string
}

func NewMockHasher() *MockHasher {
	return &MockHasher{
		HashMap: make(map[string]string),
	}
}

func (mh *MockHasher) Add(password, hash string) *MockHasher {
	mh.HashMap[password] = hash
	return mh
}

func (mh *MockHasher) Hash(password string) string {
	if hash, ok := mh.HashMap[password]; ok {
		return hash
	}
	return ""
}

func (mh *MockHasher) Verify(password, hash string) (bool, error) {
	test, ok := mh.HashMap[password]
	if !ok {
		return false, fmt.Errorf("mock password not found")
	}
	return test == hash, nil
}
