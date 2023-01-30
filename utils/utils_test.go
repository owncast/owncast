package utils

import (
	"testing"
)

func TestGetHashtagsFromText(t *testing.T) {
	text := `Some text with a #hashtag goes here.\n\n
	Another #secondhashtag, goes here.\n\n
	#thirdhashtag`

	hashtags := GetHashtagsFromText(text)

	if hashtags[0] != "#hashtag" || hashtags[1] != "#secondhashtag" || hashtags[2] != "#thirdhashtag" {
		t.Error("Incorrect hashtags fetched from text.")
	}
}

func TestPercentageUtilsTest(t *testing.T) {
	total := 42
	number := 18

	percent := IntPercentage(number, total)

	if percent != 42 {
		t.Error("Incorrect percentage calculation.")
	}
}
