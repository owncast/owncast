package models

import (
	"bytes"
	"encoding/gob"
)

// ConfigEntry is the actual object saved to the database.
// The Value is encoded using encoding/gob.
type ConfigEntry struct {
	Key   string
	Value interface{}
}

// GetStringSlice will return the value as a string slice.
func (c *ConfigEntry) GetStringSlice() ([]string, error) {
	decoder := c.GetDecoder()
	var result []string
	err := decoder.Decode(&result)
	return result, err
}

// GetStringMap will return the value as a string map.
func (c *ConfigEntry) GetStringMap() (map[string]string, error) {
	decoder := c.GetDecoder()
	var result map[string]string
	err := decoder.Decode(&result)
	return result, err
}

// GetString will return the value as a string.
func (c *ConfigEntry) GetString() (string, error) {
	decoder := c.GetDecoder()
	var result string
	err := decoder.Decode(&result)
	return result, err
}

// GetNumber will return the value as a float64.
func (c *ConfigEntry) GetNumber() (float64, error) {
	decoder := c.GetDecoder()
	var result float64
	err := decoder.Decode(&result)
	return result, err
}

// GetBool will return the value as a bool.
func (c *ConfigEntry) GetBool() (bool, error) {
	decoder := c.GetDecoder()
	var result bool
	err := decoder.Decode(&result)
	return result, err
}

// GetObject will return the value as an object.
func (c *ConfigEntry) GetObject(result interface{}) error {
	decoder := c.GetDecoder()
	err := decoder.Decode(result)
	return err
}

// GetDecoder will return a decoder for the value.
func (c *ConfigEntry) GetDecoder() *gob.Decoder {
	valueBytes := c.Value.([]byte)
	decoder := gob.NewDecoder(bytes.NewBuffer(valueBytes))
	return decoder
}
