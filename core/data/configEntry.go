package data

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

func (c *ConfigEntry) getStringSlice() ([]string, error) {
	decoder := c.getDecoder()
	var result []string
	err := decoder.Decode(&result)
	return result, err
}

func (c *ConfigEntry) getString() (string, error) {
	decoder := c.getDecoder()
	var result string
	err := decoder.Decode(&result)
	return result, err
}

func (c *ConfigEntry) getNumber() (float64, error) {
	decoder := c.getDecoder()
	var result float64
	err := decoder.Decode(&result)
	return result, err
}

func (c *ConfigEntry) getBool() (bool, error) {
	decoder := c.getDecoder()
	var result bool
	err := decoder.Decode(&result)
	return result, err
}

func (c *ConfigEntry) getObject(result interface{}) error {
	decoder := c.getDecoder()
	err := decoder.Decode(result)
	return err
}

func (c *ConfigEntry) getDecoder() *gob.Decoder {
	valueBytes := c.Value.([]byte)
	decoder := gob.NewDecoder(bytes.NewBuffer(valueBytes))
	return decoder
}
