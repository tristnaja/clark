package main

import (
	"database/sql"
	"fmt"
)

type Assistant struct {
	Name            string
	IsAssistantOn   bool
	MasterContext   string
	ImportantPeople map[string]string
	GlobalDB        *sql.DB
}

func NewAssistant(db *sql.DB) *Assistant {
	return &Assistant{
		Name:            "Clark",
		IsAssistantOn:   false,
		MasterContext:   "",
		ImportantPeople: make(map[string]string),
		GlobalDB:        db,
	}
}

func (ast *Assistant) CreateDB() error {
	query := `CREATE TABLE IF NOT EXISTS important_people (
		jid TEXT PRIMARY KEY,
		relation TEXT
	);`

	_, err := ast.GlobalDB.Exec(query)
	if err != nil {
		return fmt.Errorf("fail to create table <important_people>: %w", err)
	}

	return nil
}

func (ast *Assistant) LoadDB() error {
	var jid, relation string
	query := `SELECT jid, relations FROM important_people`

	rows, err := ast.GlobalDB.Query(query)
	if err != nil {
		return fmt.Errorf("fail to load table <important_people>: %w", err)
	}

	for rows.Next() {
		err := rows.Scan(&jid, &relation)
		if err != nil {
			return fmt.Errorf("fail to scan jid and relation: %w", err)
		}

		ast.ImportantPeople[jid] = relation
	}

	return nil
}

// TODO: AddPeople() Function
