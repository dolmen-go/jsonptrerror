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
	if err, ok := err.(*jsonptrerror.UnmarshalTypeError); ok {
		fmt.Println("Original error:", err.UnmarshalTypeError.Error())
		fmt.Println("Error location:", err.Pointer)
	}

	// Output:
	// /value: cannot unmarshal number into Go value of type bool
	// Original error: json: cannot unmarshal number into Go struct field .value of type bool
	// Error location: /value
}
