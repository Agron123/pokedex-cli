package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "mano tevas yra",
			expected: []string{"mano", "tevas", "yra"},
		},
		{
			input:    "as myliu medzius",
			expected: []string{"as", "myliu", "medzius"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Length of slices are not equal")
			continue
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Test for %s failed at index:%d, Word does not match, expected:%s, got:%s", c.input, i, expectedWord, word)
			}
		}
	}

}
