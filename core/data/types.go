package data

func (ds *Datastore) GetString(key string) (string, error) {
	configEntry, err := ds.Get(key)
	if err != nil {
		return "", err
	}
	return configEntry.getString()
}

func (ds *Datastore) SetString(key string, value string) error {
	configEntry := ConfigEntry{key, value}
	return ds.Save(configEntry)
}

func (ds *Datastore) GetNumber(key string) (float64, error) {
	configEntry, err := ds.Get(key)
	if err != nil {
		return 0, err
	}
	return configEntry.getNumber()
}

func (ds *Datastore) SetNumber(key string, value float64) error {
	configEntry := ConfigEntry{key, value}
	return ds.Save(configEntry)
}

func (ds *Datastore) GetBool(key string) (bool, error) {
	configEntry, err := ds.Get(key)
	if err != nil {
		return false, err
	}
	return configEntry.getBool()
}

func (ds *Datastore) SetBool(key string, value bool) error {
	configEntry := ConfigEntry{key, value}
	return ds.Save(configEntry)
}
