package utils

import (
	"fmt"
	"html"
	"strings"
	"testing"
)

// TestStripHTML tests the StripHTML function.
func TestStripHTML(t *testing.T) {
	requestedString := `<p><img src="img.png"/>Some text</p>`
	expectedResult := `Some text`

	result := StripHTML(requestedString)
	fmt.Println(result)

	if result != expectedResult {
		t.Errorf("Expected %s, got %s", expectedResult, result)
	}
}

// TestSafeString tests the TestSafeString function.
func TestSafeString(t *testing.T) {
	requestedString := `<p><img src="img.png"/>   Some text blah blah blah blah blah blahb albh</p>`
	expectedResult := `Some te`

	result := MakeSafeStringOfLength(requestedString, 10)
	fmt.Println(result)

	if result != expectedResult {
		t.Errorf("Expected %s, got %s", expectedResult, result)
	}
}

func TestMakeSafeStringDecodesEntitiesBeforeStrippingHTML(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: `&lt;img src=x onerror=alert(1)&gt;Alice`, want: "Alice"},
		{input: `&lt;b&gt;Alice&lt;/b&gt;`, want: "Alice"},
		{input: `Alice &amp; Bob`, want: "Alice & Bob"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			if got := MakeSafeStringOfLength(test.input, 100); got != test.want {
				t.Fatalf("MakeSafeStringOfLength(%q) = %q, want %q", test.input, got, test.want)
			}
		})
	}

	nested := `<img src=x onerror=alert(1)>`
	for range 10 {
		nested = html.EscapeString(nested)
	}
	if got := MakeSafeStringOfLength(nested, 100); strings.Contains(got, "<img") {
		t.Fatalf("deeply encoded HTML became active markup: %q", got)
	}
}

// TestSafeStringWhitespace tests the MakeSafeStringOfLength function with whitespace-only input.
func TestSafeStringWhitespace(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
		name     string
	}{
		{
			input:    "   ",
			expected: "",
			name:     "spaces only",
		},
		{
			input:    "\t\n\r",
			expected: "",
			name:     "tab, newline, carriage return",
		},
		{
			input:    " \t \n \r ",
			expected: "",
			name:     "mixed whitespace",
		},
		{
			input:    "",
			expected: "",
			name:     "empty string",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := MakeSafeStringOfLength(tc.input, 30)
			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}
