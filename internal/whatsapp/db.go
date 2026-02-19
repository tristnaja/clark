package whatsapp

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type Database struct {
	DB  *sql.DB
	Ctx context.Context
}

func InitDB() (*Database, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("fail to load .env: %w", err)
	}

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New(context.Background(), "sqlite3", "file:mystore.db?_foreign_keys=on", dbLog)

	if err != nil {
		return nil, fmt.Errorf("fail to initiate database container: %w", err)
	}

	defer container.Close()

	rawDb, err := sql.Open("sqlite3", "mystore.db")

	if err != nil {
		return nil, fmt.Errorf("fail to open database: %w", err)
	}

	if err := rawDb.Ping(); err != nil {
		return nil, fmt.Errorf("fail to ping database: %w", err)
	}

	db := &Database{
		DB: rawDb,
	}

	err = db.createDB()

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *Database) createDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	query := `CREATE TABLE IF NOT EXISTS vip (
		jid TEXT PRIMARY KEY,
		name TEXT,
		relation TEXT
	);`

	astQuery := `CREATE TABLE IF NOT EXISTS assistant_setting (
		key TEXT PRIMARY KEY,
		value TEXT
	);`

	chatQuery := `CREATE TABLE IF NOT EXISTS chat_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		jid TEXT,
		role TEXT,
		content TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.DB.ExecContext(ctx, query)

	if err != nil {
		return fmt.Errorf("fail to create table <vip>: %w", err)
	}

	_, err = db.DB.ExecContext(ctx, astQuery)

	if err != nil {
		return fmt.Errorf("fail to create table <assistant_setting>: %w", err)
	}

	_, err = db.DB.ExecContext(ctx, chatQuery)

	if err != nil {
		return fmt.Errorf("fail to create table <chat_history>: %w", err)
	}

	return nil
}

func (db *Database) SaveMessage(jid, role, content string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	saveQuery := `INSERT INTO chat_history (jid, role, content) VALUES (?, ?, ?)`
	_, err = tx.ExecContext(ctx, saveQuery, jid, role, content)
	if err != nil {
		return fmt.Errorf("fail to save message: %w", err)
	}

	cleanupQuery := `
        DELETE FROM chat_history 
        WHERE jid = ? AND id NOT IN (
            SELECT id FROM (
                SELECT id FROM chat_history 
                WHERE jid = ? 
                ORDER BY timestamp DESC 
                LIMIT 30
            ) AS temp
        )`

	_, err = tx.ExecContext(ctx, cleanupQuery, jid, jid)
	if err != nil {
		return fmt.Errorf("fail to clean up: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
