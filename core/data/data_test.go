package data

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	dbFile, err := os.CreateTemp(os.TempDir(), "owncast-test-db.db")
	if err != nil {
		panic(err)
	}

	SetupPersistence(dbFile.Name())
	m.Run()
}

func TestString(t *testing.T) {
	const testKey = "test string key"
	const testValue = "test string value"

	if err := _datastore.SetString(testKey, testValue); err != nil {
		panic(err)
	}

	// Get the config entry from the database
	stringTestResult, err := _datastore.GetString(testKey)
	if err != nil {
		panic(err)
	}

	if stringTestResult != testValue {
		t.Error("expected", testValue, "but test returned", stringTestResult)
	}
}

func TestNumber(t *testing.T) {
	const testKey = "test number key"
	const testValue = 42

	err := _datastore.SetNumber(testKey, testValue)
	if err != nil {
		panic(err)
	}

	// Get the config entry from the database
	numberTestResult, err := _datastore.GetNumber(testKey)
	if err != nil {
		panic(err)
	}
	fmt.Println(numberTestResult)

	if numberTestResult != testValue {
		t.Error("expected", testValue, "but test returned", numberTestResult)
	}
}

func TestBool(t *testing.T) {
	const testKey = "test bool key"
	const testValue = true

	err := _datastore.SetBool(testKey, testValue)
	if err != nil {
		panic(err)
	}

	// Get the config entry from the database
	numberTestResult, err := _datastore.GetBool(testKey)
	if err != nil {
		panic(err)
	}
	fmt.Println(numberTestResult)

	if numberTestResult != testValue {
		t.Error("expected", testValue, "but test returned", numberTestResult)
	}
}

func TestCustomType(t *testing.T) {
	const testKey = "test custom type key"

	// Test an example struct with a slice
	testStruct := TestStruct{
		Test:      "Test string 123 in test struct",
		TestSlice: []string{"test string 1", "test string 2"},
	}

	// Save config entry to the database
	if err := _datastore.Save(ConfigEntry{testKey, &testStruct}); err != nil {
		t.Error(err)
	}

	// Get the config entry from the database
	entryResult, err := _datastore.Get(testKey)
	if err != nil {
		t.Error(err)
	}

	// Get a typed struct out of it
	var testResult TestStruct
	if err := entryResult.getObject(&testResult); err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v", testResult)

	if testResult.TestSlice[0] != testStruct.TestSlice[0] {
		t.Error("expected", testStruct.TestSlice[0], "but test returned", testResult.TestSlice[0])
	}
}

// Custom type for testing
type TestStruct struct {
	Test            string
	TestSlice       []string
	privateProperty string
}
