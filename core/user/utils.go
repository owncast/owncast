package user

func execSQL(schemaSQL string) error {
	stmt, err := _db.Prepare(schemaSQL)
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
