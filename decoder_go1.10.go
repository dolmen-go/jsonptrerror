//+build go1.10

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
