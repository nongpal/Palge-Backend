package data

import "time"

type Token struct {
	Plaintext string
	Hash      []byte
	UserID    int64
	Expiry    time.Time
	Scope     string
}
