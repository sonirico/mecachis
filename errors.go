package mecachis

import "fmt"

type ErrDuplicatedKey struct {
	key string
}

func NewDuplicatedKeyError(key string) *ErrDuplicatedKey {
	return &ErrDuplicatedKey{key: key}
}

func (e *ErrDuplicatedKey) Error() string {
	return fmt.Sprintf("'%s' is already in the cache", e.key)
}
