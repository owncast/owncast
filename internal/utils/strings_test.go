package utils

import (
	"fmt"
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
