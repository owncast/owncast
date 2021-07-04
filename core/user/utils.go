package user

func execSQL(schemaSQL string) error {
	_datastore.DbLock.Lock()
	defer _datastore.DbLock.Unlock()

	stmt, err := _datastore.DB.Prepare(schemaSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}
