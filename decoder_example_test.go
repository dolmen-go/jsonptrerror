// The original error message returned by stdlib changed with go1.8.
// We only test the latest release.
//
//go:build go1.8 || forcego1.8
// +build go1.8 forcego1.8

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
		//fmt.Println("Original error:", err.UnmarshalTypeError.Error())
		fmt.Println("Error location:", err.Pointer)
	}

	// Output:
	// /value: cannot unmarshal number into Go value of type bool
	// Error location: /value
}
