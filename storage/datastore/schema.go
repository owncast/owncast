package datastore

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

// This is a central point for creating all the tables required for the application.
func (ds *Datastore) createTables() {
	createDatastoreTable(ds.DB)
	createUsersTable(ds.DB)
	createAccessTokenTable(ds.DB)
	createWebhooksTable(ds.DB)
	createFederationFollowersTable(ds.DB)
	createFederationOutboxTable(ds.DB)
	createNotificationsTable(ds.DB)
	CreateBanIPTable(ds.DB)
	CreateMessagesTable(ds.DB)
}

func createDatastoreTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS datastore (
		"key" string NOT NULL PRIMARY KEY,
		"value" BLOB,
		"timestamp" DATE DEFAULT CURRENT_TIMESTAMP NOT NULL
	);`

	MustExec(createTableSQL, db)
}

func createWebhooksTable(db *sql.DB) {
	log.Traceln("Creating webhooks table...")

	createTableSQL := `CREATE TABLE IF NOT EXISTS webhooks (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"url" string NOT NULL,
		"events" TEXT NOT NULL,
		"timestamp" DATETIME DEFAULT CURRENT_TIMESTAMP,
		"last_used" DATETIME
	);`

	stmt, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	if _, err = stmt.Exec(); err != nil {
		log.Warnln(err)
	}
}

func createUsersTable(db *sql.DB) {
	log.Traceln("Creating users table...")

	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"id" TEXT,
		"display_name" TEXT NOT NULL,
		"display_color" NUMBER NOT NULL,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		"disabled_at" TIMESTAMP,
		"previous_names" TEXT DEFAULT '',
		"namechanged_at" TIMESTAMP,
    "authenticated_at" TIMESTAMP,
		"scopes" TEXT,
		"type" TEXT DEFAULT 'STANDARD',
		"last_used" DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
	);`

	MustExec(createTableSQL, db)
	MustExec(`CREATE INDEX IF NOT EXISTS idx_user_id ON users (id);`, db)
	MustExec(`CREATE INDEX IF NOT EXISTS idx_user_id_disabled ON users (id, disabled_at);`, db)
	MustExec(`CREATE INDEX IF NOT EXISTS idx_user_disabled_at ON users (disabled_at);`, db)
}

func createAccessTokenTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS user_access_tokens (
    "token" TEXT NOT NULL PRIMARY KEY,
    "user_id" TEXT NOT NULL,
    "timestamp" DATE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id)
  );`

	stmt, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		log.Warnln(err)
	}
}

func createFederationFollowersTable(db *sql.DB) {
	log.Traceln("Creating federation followers table...")

	createTableSQL := `CREATE TABLE IF NOT EXISTS ap_followers (
		"iri" TEXT NOT NULL,
		"inbox" TEXT NOT NULL,
		"name" TEXT,
		"username" TEXT NOT NULL,
		"image" TEXT,
    "request" TEXT NOT NULL,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		"approved_at" TIMESTAMP,
    "disabled_at" TIMESTAMP,
    "request_object" BLOB,
		PRIMARY KEY (iri));`
	MustExec(createTableSQL, db)
	MustExec(`CREATE INDEX IF NOT EXISTS idx_iri ON ap_followers (iri);`, db)
	MustExec(`CREATE INDEX IF NOT EXISTS idx_approved_at ON ap_followers (approved_at);`, db)
}

func createFederationOutboxTable(db *sql.DB) {
	log.Traceln("Creating federation outbox table...")
	createTableSQL := `CREATE TABLE IF NOT EXISTS ap_outbox (
		"iri" TEXT NOT NULL,
		"value" BLOB,
		"type" TEXT NOT NULL,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "live_notification" BOOLEAN DEFAULT FALSE,
		PRIMARY KEY (iri));`

	MustExec(createTableSQL, db)
	MustExec(`CREATE INDEX IF NOT EXISTS idx_iri ON ap_outbox (iri);`, db)
	MustExec(`CREATE INDEX IF NOT EXISTS idx_type ON ap_outbox (type);`, db)
	MustExec(`CREATE INDEX IF NOT EXISTS idx_live_notification ON ap_outbox (live_notification);`, db)
}

// createFederatedActivitiesTable will create the accepted
// activities table if needed.
func createFederatedActivitiesTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS ap_accepted_activities (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"iri" TEXT NOT NULL,
    "actor" TEXT NOT NULL,
    "type" TEXT NOT NULL,
		"timestamp" TIMESTAMP NOT NULL
	);`

	MustExec(createTableSQL, db)
	MustExec(`CREATE INDEX IF NOT EXISTS idx_iri_actor_index ON ap_accepted_activities (iri,actor);`, db)
}

// CreateMessagesTable will create the chat messages table if needed.
func CreateMessagesTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS messages (
		"id" string NOT NULL,
		"user_id" TEXT,
		"body" TEXT,
		"eventType" TEXT,
		"hidden_at" DATETIME,
		"timestamp" DATETIME,
		"title" TEXT,
		"subtitle" TEXT,
		"image" TEXT,
		"link" TEXT,
		PRIMARY KEY (id)
	);`
	MustExec(createTableSQL, db)

	// Create indexes
	MustExec(`CREATE INDEX IF NOT EXISTS user_id_hidden_at_timestamp ON messages (id, user_id, hidden_at, timestamp);`, db)
	MustExec(`CREATE INDEX IF NOT EXISTS idx_id ON messages (id);`, db)
	MustExec(`CREATE INDEX IF NOT EXISTS idx_user_id ON messages (user_id);`, db)
	MustExec(`CREATE INDEX IF NOT EXISTS idx_hidden_at ON messages (hidden_at);`, db)
	MustExec(`CREATE INDEX IF NOT EXISTS idx_timestamp ON messages (timestamp);`, db)
	MustExec(`CREATE INDEX IF NOT EXISTS idx_messages_hidden_at_timestamp on messages(hidden_at, timestamp);`, db)
}

// CreateBanIPTable will create the IP ban table if needed.
func CreateBanIPTable(db *sql.DB) {
	createTableSQL := `  CREATE TABLE IF NOT EXISTS ip_bans (
    "ip_address" TEXT NOT NULL PRIMARY KEY,
    "notes" TEXT,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );`

	stmt, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal("error creating ip ban table", err)
	}
	defer stmt.Close()
	if _, err := stmt.Exec(); err != nil {
		log.Fatal("error creating ip ban table", err)
	}
}

func createNotificationsTable(db *sql.DB) {
	log.Traceln("Creating federation followers table...")

	createTableSQL := `CREATE TABLE IF NOT EXISTS notifications (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"channel" TEXT NOT NULL,
		"destination" TEXT NOT NULL,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`

	MustExec(createTableSQL, db)
	MustExec(`CREATE INDEX IF NOT EXISTS idx_channel ON notifications (channel);`, db)
}
