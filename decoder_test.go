package jsonptrerror_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/dolmen-go/jsonptrerror"
)

func TestDecoder(t *testing.T) {
	for _, test := range []struct {
		in    string
		value interface{}
		ptr   interface{}
	}{
		{`"a"`, new(string), nil},
		{`true`, new(bool), nil},
		{`null`, new(interface{}), nil},
		{`[]`, new([]interface{}), nil},
		{`{}`, new(map[string]interface{}), nil},
		{`"a"`, new(bool), ""},
		{`[]`, new(map[string]interface{}), nil},
		{`["a"]`, &([]string{}), nil},
		{`["a"]`, &([]int{}), "/0"},
		{`[1]`, &([]int{}), nil},
		{`[1]`, &([]string{}), "/0"},
		{`[1]`, &([]bool{}), "/0"},
		{`[true]`, &([]int{}), "/0"},
		{`[true]`, &([]string{}), "/0"},
		{`[1,true]`, &([]int{}), "/1"},
		{`["a",true]`, &([]string{}), "/1"},
		{`[1,true,1]`, &([]int{}), "/1"},
		{`{}`, new([]interface{}), nil},
		{`{"a":1}`, new(map[string]interface{}), nil},
		{`{"a":1}`, new(map[string]string), "/a"},
		{`{"~":1}`, new(map[string]string), "/~0"},
		{`{"/":1}`, new(map[string]string), "/~1"},
		{`{"/":[1]}`, new(map[string][]string), "/~1/0"},
		{` {  "/" : [ 1]}`, new(map[string][]string), "/~1/0"},
		{` {  "\u002f" : [ 1]}`, new(map[string][]string), "/~1/0"},
		{`[{},{"a":1}]`, new([]map[string]int), nil},
		{`[{},{"a":1}]`, new([]map[string]bool), "/1/a"},
		// TODO structs
	} {
		t.Logf("%s -> %T", test.in, test.value)
		//var v1, v2 interface{}
		err1 := json.NewDecoder(bytes.NewBufferString(test.in)).Decode(test.value)
		err2 := jsonptrerror.NewDecoder(bytes.NewBufferString(test.in)).Decode(test.value)
		if (err1 == nil) != (err2 == nil) {
			t.Errorf("err = %q, want: %q", err2, err1)
		} else if test.ptr != nil {
			if err2, ok := err2.(*jsonptrerror.UnmarshalTypeError); ok {
				if err2.Pointer.String() != test.ptr.(string) {
					t.Errorf("ptr = %q, want: %q", err2.Pointer, test.ptr)
				}
			} else {
				t.Errorf("err = %q, want jsponptrerror.UnmarshalTypeError", err2)
			}
		}
	}
}
