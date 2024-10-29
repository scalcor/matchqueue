package matchqueue

import "errors"

var (
	ErrNotInitialized  = errors.New("not initialized")
	ErrNotEnoughPlayer = errors.New("not enough player")
)
