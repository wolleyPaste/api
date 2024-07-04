package db

import (
	"database/sql"
	"fmt"
	"wolley-api/src/common"

	_ "github.com/lib/pq"
)

var db *sql.DB

// InitPostgres initializes the PostgreSQL connection
func InitPostgres(postgreSQLDB common.PostgreSQL) error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		postgreSQLDB.Host, postgreSQLDB.Port, postgreSQLDB.Username, postgreSQLDB.Password, postgreSQLDB.DBName)
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	return db.Ping()
}

func ClosePostgres() {
	db.Close()
}

// Checks if a paste with a specefic ID exists
func GetPaste(ID string) (common.Paste, bool) {
	var text string
	var title string
	var password string
	var creationString string
	var expirationString string

	err := db.QueryRow(`
            SELECT text, title, password, creation, expiration
            FROM pastes 
            WHERE ID = $1
    `, ID).Scan(&text, &title, &password, &creationString, &expirationString)

	if err != nil {
		return common.Paste{}, false
	}

	creation := common.StringToTime(creationString)
	expiration := common.StringToTime(expirationString)

	paste := common.Paste{
		ID:           ID,
		Text:         text,
		Title:        title,
		Password:     password,
		CreationDate: *creation,
	}

	if expiration != nil {
		paste.ExpirationDate = *expiration
	}

	return paste, true
}

// Checks if a paste with a specefic ID exists
func CheckPasteExistence(ID string) bool {
	var exists bool
	err := db.QueryRow(`
        SELECT EXISTS (
            SELECT 1 
            FROM pastes 
            WHERE ID = $1
        )
    `, ID).Scan(&exists)

	if err != nil {
		return false
	}

	return exists
}
