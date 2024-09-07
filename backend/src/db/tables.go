package db

// This only initializes tables, it does not handle migrations
func (db *Database) InitializeTables() error {
	var err error

	err = db.initalizeVersionTable()
	if err != nil {
		return err
	}

	err = db.initializeUserTable()
	if err != nil {
		return err
	}

	err = db.initializeSourcesTable()
	if err != nil {
		return err
	}

	return nil
}
