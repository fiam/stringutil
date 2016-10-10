package stringutil

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type splitCase struct {
	s      string
	sep    string
	result []string
}

func sepRepr(sep string) string {
	switch sep {
	case "":
		return "SPACE"
	}
	return fmt.Sprintf("%q", sep)
}

func resultRepr(res []string) string {
	values := make([]string, len(res))
	for ii, v := range res {
		values[ii] = fmt.Sprintf("%q", v)
	}
	return fmt.Sprintf("[%s] (%d values)", strings.Join(values, ", "), len(values))
}

func TestSplitFields(t *testing.T) {
	cases := []splitCase{
		{"The, quick, brown", ",", []string{"The", "quick", "brown"}},
		{"'fo\"x', 'jum,ps', \"ov',er\"", ",", []string{"fo\"x", "jum,ps", "ov',er"}},
		{"  a\tb\r\nc ", "", []string{"a", "b", "c"}},
		{"''  a\tb\r\nc ", "", []string{"", "a", "b", "c"}},
		{"''  a\tb\r\nc ''   ", "", []string{"", "a", "b", "c", ""}},
		{"a = b", "=", []string{"a", "b"}},
		{" a = b ", "=", []string{"a", "b"}},
		{"'a ' = b", "=", []string{"a ", "b"}},
		{"'a ' = ' b '", "=", []string{"a ", " b "}},
		{"' a' = b", "=", []string{" a", "b"}},
		{"' a ' = b", "=", []string{" a ", "b"}},
		{"'a' = 'b' ", "=", []string{"a", "b"}},
	}
	for _, v := range cases {
		fields, err := SplitFields(v.s, v.sep)
		if err != nil {
			t.Errorf("error splitting %q with sep %s: %s", v.s, sepRepr(v.sep), err)
			continue
		}
		if !reflect.DeepEqual(fields, v.result) {
			t.Errorf("error splitting %q with sep %s. wanted %v, got %v", v.s, sepRepr(v.sep), resultRepr(v.result), resultRepr(fields))
		}
	}
}

func TestKeepQuotes(t *testing.T) {
	fields, err := SplitFieldsOptions("\"The\", 'quick', brown", ",", &SplitOptions{KeepQuotes: true})
	if err != nil {
		t.Fatal(err)
	}
	exp := []string{"\"The\"", "'quick'", "brown"}
	if !reflect.DeepEqual(fields, exp) {
		t.Errorf("error splitting keeping quotes - want %v, got %v", exp, fields)
	}
}

func TestSplitLines(t *testing.T) {
	tests := map[string][]string{
		"a\nb\nc":     []string{"a", "b", "c"},
		"a\r\nb\nc":   []string{"a", "b", "c"},
		"a\r\nb\\\nc": []string{"a", "bc"},
	}
	for k, v := range tests {
		lines := SplitLines(k)
		if !reflect.DeepEqual(lines, v) {
			t.Errorf("expecting %v when splitting lines from %q, got %v instead", v, k, lines)
		}
	}
}

func TestSplitCommonPrefix(t *testing.T) {
	type splitTestCase struct {
		input  []string
		prefix string
		output []string
	}
	cases := []splitTestCase{
		{[]string{"a", "b", "c"}, "", []string{"a", "b", "c"}},
		{[]string{"aa", "ab", "ac"}, "a", []string{"a", "b", "c"}},
		{[]string{"aaaa", "aaaa", ""}, "", []string{"aaaa", "aaaa", ""}},
		{[]string{"aaa", "ab", "ac"}, "a", []string{"aa", "b", "c"}},
		{[]string{"", "", ""}, "", []string{"", "", ""}},
		{[]string{"go"}, "go", []string{""}},
		{nil, "", nil},
	}
	for _, v := range cases {
		prefix, output := SplitCommonPrefix(v.input)
		if prefix != v.prefix {
			t.Errorf("expecting prefix %q from input %v, got %q instead", v.prefix, v.input, prefix)
		}
		if !reflect.DeepEqual(output, v.output) {
			t.Errorf("expecting output %v from input %v, got %v instead", v.output, v.input, output)
		}
	}
}
