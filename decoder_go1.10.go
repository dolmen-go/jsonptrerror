//go:build go1.10 && !go1.14
// +build go1.10,!go1.14

package jsonptrerror

import (
	"encoding/json"
	"io"
)

type decoder interface {
	Decode(interface{}) error
	UseNumber()
	DisallowUnknownFields()
	More() bool
	Buffered() io.Reader
	Token() (json.Token, error)
}
