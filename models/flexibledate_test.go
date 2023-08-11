package models

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"
)

func TestFlexibleDateParsing(t *testing.T) {
	type testJson struct {
		Testdate FlexibleDate `json:"testdate"`
	}

	nullTime := sql.NullTime{Time: time.Unix(1591614434, 0), Valid: true}
	testNullTime, err := FlexibleDateParse(nullTime)
	if err != nil {
		t.Error(err)
	}

	if testNullTime.Unix() != nullTime.Time.Unix() {
		t.Errorf("Expected %d but got %d", nullTime.Time.Unix(), testNullTime.Unix())
	}

	testStrings := map[string]time.Time{
		"2023-08-10 17:40:15.376736475-07:00": time.Unix(1691714415, 0),
	}

	for testString, expectedTime := range testStrings {
		testJsonString := `{"testdate":"` + testString + `"}`
		response := testJson{}

		err := json.Unmarshal([]byte(testJsonString), &response)
		if err != nil {
			t.Error(err)
		}

		if response.Testdate.Time.Unix() != expectedTime.Unix() {
			t.Errorf("Expected %d but got %d", expectedTime.Unix(), response.Testdate.Time.Unix())
		}
	}
}
