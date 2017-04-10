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

func NewDecoder(r io.Reader) *Decoder {
	var d Decoder
	d.decoder = json.NewDecoder(io.TeeReader(r, &d.input))
	return &d
}

func (d *Decoder) UseNumber() {
	d.decoder.UseNumber()
}

const (
	tokenTopValue = iota
	tokenArrayValue
	tokenArrayComma
	tokenObjectKey
	tokenObjectColon
	tokenObjectValue
	tokenObjectComma
)

func (d *Decoder) Decode(v interface{}) error {
	if d.err != nil {
		return d.err
	}
	d.err = d.decoder.Decode(v)
	if err, ok := d.err.(*json.UnmarshalTypeError); ok {
		offset := int(err.Offset)
		input := d.input.Bytes()
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
				ptr := jsonptr.Pointer{}
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
				d.err = &UnmarshalTypeError{*err, ptr}
				break
			}
			switch input[i] {
			case '{':
				elemStack = append(elemStack, elem{container: '{'})
				expectKey = true
			case '[':
				elemStack = append(elemStack, elem{container: '[', index: 0})
			case '}', ']':
				elemStack = elemStack[:len(elemStack)-1]
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
	}
	if d.err != nil {
		d.decoder = nil
		d.input = bytes.Buffer{}
	}
	return d.err
}
