package jsonptrerror_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sort"
	"strings"
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

func listTypes(num int, at func(int) reflect.Type) []string {
	if num == 0 {
		return nil
	}
	res := make([]string, 0, num)
	for i := 0; i < num; i++ {
		res = append(res, at(i).String())
	}
	return res
}

func listMethods(p interface{}) []string {
	t := reflect.ValueOf(p).Type()
	meths := make([]string, 0, t.NumMethod())
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		// Skip unexported methods
		if len(m.PkgPath) != 0 {
			continue
		}
		mt := m.Type
		in := listTypes(mt.NumIn()-1, func(i int) reflect.Type { return mt.In(i + 1) })
		out := listTypes(mt.NumOut(), mt.Out)
		sig := bytes.NewBufferString(m.Name)
		sig.WriteByte('(')
		if in != nil {
			sig.WriteString(strings.Join(in, ", "))
		}
		sig.WriteByte(')')
		switch len(out) {
		case 0:
		case 1:
			sig.WriteByte(' ')
			sig.WriteString(out[0])
		default:
			sig.WriteString(" (")
			sig.WriteString(strings.Join(out, ", "))
			sig.WriteByte(')')
		}
		meths = append(meths, sig.String())
	}
	sort.Strings(meths)
	return meths
}

func TestInterface(t *testing.T) {
	m1 := listMethods(&json.Decoder{})
	m2 := listMethods(&jsonptrerror.Decoder{})
	if len(m1) == 0 {
		t.Fatal("bug in listMethods")
	}
	if len(m1) != len(m2) {
		t.Fatalf("%v != %v", m1, m2)
	}
	for i := range m1 {
		if m1[i] != m2[i] {
			t.Fatalf("decoder interface doesn't match json.Decoder: %#v != %#v\n"+
				"Check: grep 'pkg encoding/json, method (\\*Decoder)' $GOROOT/api/go*.txt",
				m1[i], m2[i])
		}
	}
}
