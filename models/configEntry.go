package models

import (
	"bytes"
	"encoding/gob"
)

// ConfigEntry is the actual object saved to the database.
// The Value is encoded using encoding/gob.
type ConfigEntry struct {
	Value interface{}
	Key   string
}

func (c *ConfigEntry) GetStringSlice() ([]string, error) {
	decoder := c.GetDecoder()
	var result []string
	err := decoder.Decode(&result)
	return result, err
}

func (c *ConfigEntry) GetStringMap() (map[string]string, error) {
	decoder := c.GetDecoder()
	var result map[string]string
	err := decoder.Decode(&result)
	return result, err
}

func (c *ConfigEntry) GetString() (string, error) {
	decoder := c.GetDecoder()
	var result string
	err := decoder.Decode(&result)
	return result, err
}

func (c *ConfigEntry) GetNumber() (float64, error) {
	decoder := c.GetDecoder()
	var result float64
	err := decoder.Decode(&result)
	return result, err
}

func (c *ConfigEntry) GetBool() (bool, error) {
	decoder := c.GetDecoder()
	var result bool
	err := decoder.Decode(&result)
	return result, err
}

func (c *ConfigEntry) GetObject(result interface{}) error {
	decoder := c.GetDecoder()
	err := decoder.Decode(result)
	return err
}

func (c *ConfigEntry) GetDecoder() *gob.Decoder {
	valueBytes := c.Value.([]byte)
	decoder := gob.NewDecoder(bytes.NewBuffer(valueBytes))
	return decoder
}
