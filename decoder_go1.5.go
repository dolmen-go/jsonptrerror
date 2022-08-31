//go:build go1.5 && !go1.10
// +build go1.5,!go1.10

package jsonptrerror

import (
	"encoding/json"
	"io"
)

type decoder interface {
	Decode(interface{}) error
	UseNumber()
	More() bool
	Buffered() io.Reader
	Token() (json.Token, error)
}
