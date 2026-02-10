package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
)

type Assistant struct {
	Name          string
	Status        bool
	MasterContext string
	VIP           map[string]string
	GlobalDB      *sql.DB
}

func NewAssistant(db *sql.DB) (*Assistant, error) {
	ast := &Assistant{
		Name:          "",
		Status:        false,
		MasterContext: "",
		VIP:           make(map[string]string),
		GlobalDB:      db,
	}

	err := ast.CreateDB()
	if err != nil {
		return nil, err
	}
	fmt.Println("| DB detected")

	err = ast.LoadVIP()
	if err != nil {
		return nil, err
	}
	fmt.Println("| VIP detected")

	ast.LoadAssistant()
	if ast.Name != "NULL" {
		var name string
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("You have no active assistant yet")
		fmt.Println("To initiate an assistant, it's easy, you just have name it.")

		fmt.Print("Name your assistant: ")
		name, _ = reader.ReadString('\n')

		ast.Name = name
		fmt.Printf("Assistant: %v", ast.Name)
	}
	fmt.Println("| Assistant Active")

	return ast, nil
}

func (ast *Assistant) CreateDB() error {
	query := `CREATE TABLE IF NOT EXISTS vip (
		jid TEXT PRIMARY KEY,
		relation TEXT
	);`

	astQuery := `CREATE TABLE IF NOT EXISTS assistant_setting (
		key TEXT PRIMARY KEY,
		value TEXT,
	);`

	_, err := ast.GlobalDB.Exec(query)
	if err != nil {
		return fmt.Errorf("fail to create table <vip>: %w", err)
	}

	_, err = ast.GlobalDB.Exec(astQuery)
	if err != nil {
		return fmt.Errorf("fail to create table <assistant_setting>: %w", err)
	}

	return nil
}

func (ast *Assistant) LoadVIP() error {
	var jid, relation string
	query := `SELECT jid, relations FROM vip`

	rows, err := ast.GlobalDB.Query(query)
	if err != nil {
		return fmt.Errorf("fail to load table <vip>: %w", err)
	}

	for rows.Next() {
		err := rows.Scan(&jid, &relation)
		if err != nil {
			return fmt.Errorf("fail to scan jid and relation: %w", err)
		}

		ast.VIP[jid] = relation
	}

	return nil
}

func (ast *Assistant) LoadAssistant() error {
	var name, status, context string
	ast.GlobalDB.QueryRow("SELECT value FROM assistant_setting WHERE key = 'name'").Scan(&name)
	ast.GlobalDB.QueryRow("SELECT value FROM assistant_setting WHERE key = 'status'").Scan(&status)
	ast.GlobalDB.QueryRow("SELECT value FROM assistant_setting WHERE key = 'context'").Scan(&context)

	if name == "" {
		ast.Name = "NULL"
	} else {
		ast.Name = name
	}

	if status == "" {
		ast.Status = false
	} else {
		ast.Status = true
	}

	if context == "" {
		ast.MasterContext = ""
	} else {
		ast.MasterContext = context
	}

	return nil
}

func (ast *Assistant) AddVIP(jid, relation string) error {
	query := "INSERT OR REPLACE INTO important_people (jid, relation) VALUES (?, ?)"

	_, err := ast.GlobalDB.Exec(query, jid, relation)
	if err != nil {
		return fmt.Errorf("fail to add new vip: %w", err)
	}

	err = ast.LoadVIP()
	if err != nil {
		return fmt.Errorf("fail to load new vip: %w", err)
	}

	return nil
}
