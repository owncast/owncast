package data

import (
	"testing"
)

func TestCachedString(t *testing.T) {
	const testKey = "test string key"
	const testValue = "test string value"

	_datastore.SetCachedValue(testKey, []byte(testValue))

	// Get the config entry from the database
	stringTestResult, err := _datastore.GetCachedValue(testKey)
	if err != nil {
		panic(err)
	}

	if string(stringTestResult) != testValue {
		t.Error("expected", testValue, "but test returned", stringTestResult)
	}
}
