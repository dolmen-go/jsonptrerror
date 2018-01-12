package jsonptrerror_test

import (
	"fmt"
	"strings"

	"github.com/dolmen-go/jsonptrerror"
)

func ExampleDecoder() {
	decoder := jsonptrerror.NewDecoder(strings.NewReader(
		`{"key": "x", "value": 5}`,
	))
	var out struct {
		Key   string `json:"key"`
		Value bool   `json:"value"`
	}
	err := decoder.Decode(&out)
	fmt.Println(err)

	// Output:
	// /value: cannot unmarshal number into Go value of type bool
}
