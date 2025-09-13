package storage

import (
    "database/sql"
    "io/ioutil"
    log "github.com/sirupsen/logrus"
    "path/filepath"
    _ "github.com/lib/pq"

)

// NewDB создаёт подключение к PostgreSQL
func NewDB(dataSource string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}



func InitDB(db *sql.DB) {
    var exists bool
    err := db.QueryRow(`SELECT EXISTS (
        SELECT FROM information_schema.tables 
        WHERE table_schema = 'public' AND table_name = 'tasks'
    );`).Scan(&exists)
    if err != nil {
        log.Fatal("Failed to check if table exists:", err)
    }

    if !exists {
        path := filepath.Join("schema", "tasks.sql")
        sqlBytes, err := ioutil.ReadFile(path)
        if err != nil {
            log.Fatal("Failed to read schema file:", err)
        }

        _, err = db.Exec(string(sqlBytes))
        if err != nil {
            log.Fatal("Failed to execute schema SQL:", err)
        }
    }
}
