package legacymigrations

import (
	"database/sql"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/utils"
)

func CreateNotificationsTable(db *sql.DB) {
	log.Traceln("Creating federation followers table...")

	createTableSQL := `CREATE TABLE IF NOT EXISTS notifications (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"channel" TEXT NOT NULL,
		"destination" TEXT NOT NULL,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`

	utils.MustExec(createTableSQL, db)
	utils.MustExec(`CREATE INDEX IF NOT EXISTS idx_channel ON notifications (channel);`, db)
}
