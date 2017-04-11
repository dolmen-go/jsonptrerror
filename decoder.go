/*
Package jsonptrerror extends encoding/json.Decoder to return unmarshal errors
located with JSON Pointer (RFC 6091).

The current implementation keeps a duplicate copy of the JSON document in memory.
*/
package jsonptrerror

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/dolmen-go/jsonptr"
)

// UnmarshalTypeError is an extension of encoding/json.UnmarshalTypeError
// that also includes the error location as a JSON Pointer (RFC 6901)
type UnmarshalTypeError struct {
	json.UnmarshalTypeError
	Pointer jsonptr.Pointer
}

func (e UnmarshalTypeError) Error() string {
	return e.Pointer.String() + ": cannot unmarshal " + e.Value + " into Go value of type " + e.Type.String()
}

// Unmarshal parses the JSON-encoded data and stores the result in the value pointed to by v.
//
// json.UnmarshalTypeError is translated to the extended jsonptrerror.UnmarshalTypeError
func Unmarshal(document []byte, v interface{}) error {
	err := json.Unmarshal(document, v)
	return translateError(document, err)
}

type decoder interface {
	Decode(interface{}) error
	UseNumber()
}

// Decoder is the same as encoding/json.Decoder, except Decode returns
// our UnmarshalTypeError (providing a JSON Pointer) instead of encoding/json.UnmarshalTypeError.
type Decoder struct {
	decoder decoder
	input   bytes.Buffer
	err     error
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	var d Decoder
	d.decoder = json.NewDecoder(io.TeeReader(r, &d.input))
	return &d
}

func (d *Decoder) UseNumber() {
	d.decoder.UseNumber()
}

func (d *Decoder) Decode(v interface{}) error {
	if d.err != nil {
		return d.err
	}
	d.err = d.decoder.Decode(v)
	d.err = translateError(d.input.Bytes(), d.err)
	if d.err != nil {
		d.decoder = nil
		d.input = bytes.Buffer{}
	}
	return d.err
}

func translateError(document []byte, err error) error {
	if e, ok := err.(*json.UnmarshalTypeError); ok && e != nil {
		err = &UnmarshalTypeError{*e, pointerAtOffset(document, int(e.Offset))}
	}
	return err
}

// pointerAtOffset extracts the JSON Pointer at the start of a value in a *valid* JSON document
func pointerAtOffset(input []byte, offset int) jsonptr.Pointer {
	var ptr jsonptr.Pointer
	i := 0
	type elem struct {
		container byte
		property  []byte
		index     int
	}
	var elemStack []elem
	var expectKey bool
	for {
		if i == offset {
			for _, e := range elemStack {
				switch e.container {
				case '{':
					var name string
					json.Unmarshal(e.property, &name)
					ptr.Property(name)
				case '[':
					ptr.Index(e.index)
				}
			}
			break
		}
		switch input[i] {
		case '{':
			// push state of the new current object on the stack
			elemStack = append(elemStack, elem{container: '{'})
			expectKey = true
		case '[':
			// push state of the new current array on the stack
			elemStack = append(elemStack, elem{container: '[', index: 0})
		case '}', ']':
			elemStack = elemStack[:len(elemStack)-1] // pop
		case '"':
			j := i
		str:
			for {
				i++
				switch input[i] {
				case '\\':
					i++
				case '"':
					break str
				}
			}
			if expectKey {
				elemStack[len(elemStack)-1].property = input[j : i+1]
				expectKey = false
			}
		case ',':
			if elemStack[len(elemStack)-1].container == '{' {
				expectKey = true
			} else {
				elemStack[len(elemStack)-1].index++
			}
		}
		i++
	}
	return ptr
}
