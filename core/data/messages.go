package data

// GetMessagesCount will return the number of messages in the database.
func GetMessagesCount() int64 {
	query := `SELECT COUNT(*) FROM messages`
	rows, err := _db.Query(query)
	if err != nil || rows.Err() != nil {
		return 0
	}
	defer rows.Close()
	var count int64
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0
		}
	}
	return count
}
