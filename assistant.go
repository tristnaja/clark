package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/eduardolat/openroutergo"
)

type Assistant struct {
	Name          string
	Status        bool
	MasterContext string
	VIP           map[string]string
	GlobalDB      *sql.DB
	ORClient      *openroutergo.Client
	Model         string
}

func NewAssistant(db *sql.DB, apikey string) (*Assistant, error) {
	client, err := openroutergo.NewClient().WithAPIKey(apikey).Create()

	if err != nil {
		return nil, fmt.Errorf("fail to initiate AI client: %w", err)
	}

	ast := &Assistant{
		Name:          "",
		Status:        false,
		MasterContext: "",
		VIP:           make(map[string]string),
		GlobalDB:      db,
		ORClient:      client,
		Model:         "arcee-ai/trinity-large-preview:free",
	}

	err = ast.createDB()
	if err != nil {
		return nil, err
	}
	fmt.Println("| DB detected")

	err = ast.loadVIP()
	if err != nil {
		return nil, err
	}
	fmt.Println("| VIP detected")

	ast.loadAssistant()
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

func (ast *Assistant) createDB() error {
	query := `CREATE TABLE IF NOT EXISTS vip (
		jid TEXT PRIMARY KEY,
		relation TEXT
	);`

	astQuery := `CREATE TABLE IF NOT EXISTS assistant_setting (
		key TEXT PRIMARY KEY,
		value TEXT,
	);`

	chatQuery := `CREATE TABLE IF NOT EXISTS chat_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		jid TEXT,
		role TEXT,
		content TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := ast.GlobalDB.Exec(query)
	if err != nil {
		return fmt.Errorf("fail to create table <vip>: %w", err)
	}

	_, err = ast.GlobalDB.Exec(astQuery)
	if err != nil {
		return fmt.Errorf("fail to create table <assistant_setting>: %w", err)
	}

	_, err = ast.GlobalDB.Exec(chatQuery)
	if err != nil {
		return fmt.Errorf("fail to create table <chat_history>: %w", err)
	}

	return nil
}

func (ast *Assistant) loadVIP() error {
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

func (ast *Assistant) loadAssistant() error {
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

	err = ast.loadVIP()
	if err != nil {
		return fmt.Errorf("fail to load new vip: %w", err)
	}

	return nil
}

func (ast *Assistant) SaveMessage(jid, role, content string) {
	saveQuery := `INSERT INTO chat_history (jid, role, content) VALUES (?, ?, ?)`
	_, err := ast.GlobalDB.Exec(saveQuery, jid, role, content)
	if err != nil {
		fmt.Println("fail to save message:", err)
	}

	cleanupQuery := `
        DELETE FROM chat_history 
        WHERE jid = ? AND id NOT IN (
            SELECT id FROM chat_history 
            WHERE jid = ? 
            ORDER BY timestamp DESC 
            LIMIT 30
        )`
	_, err = ast.GlobalDB.Exec(cleanupQuery, jid, jid)
	if err != nil {
		fmt.Println("fail to clean up:", err)
	}
}

func (ast *Assistant) GetHistory(jid string) []openroutergo.ChatCompletionMessage {
	getHistoryQuery := "SELECT role, content FROM chat_history WHERE jid = ? ORDER BY timestamp ASC"
	rows, _ := ast.GlobalDB.Query(getHistoryQuery, jid)
	defer rows.Close()

	var history []openroutergo.ChatCompletionMessage
	for rows.Next() {
		var dbRole, content string
		rows.Scan(&dbRole, &content)

		switch dbRole {
		case "user":
			role := openroutergo.RoleUser
			history = append(history, openroutergo.ChatCompletionMessage{
				Role:    role,
				Content: content,
			})
		case "assistant":
			role := openroutergo.RoleAssistant
			history = append(history, openroutergo.ChatCompletionMessage{
				Role:    role,
				Content: content,
			})
		}
	}
	return history
}

func (ast *Assistant) GetAIResponse(senderJid, userMsg string) string {
	relation := ast.VIP[senderJid]
	promptTemplate := `
	Act as %s, the dedicated and highly professional AI Butler for Tristan Al Harrish Basori.

	### CORE IDENTITY & ETIQUETTE
	* **Persona:** You are a refined, humble, and impeccably polite servant. Your tone is formal yet warm—think of a high-end concierge.
	* **Primary Duty:** You are currently dispatched because your Master, Tristan, is %s.
	* **Restraints:** Maintain strict professionalism. No profanity or NSFW content.

	### INTERPERSONAL GUIDELINES
	* **VIPs:** Treat these individuals with the highest reverence: %s.
	* **The Visitor:** You are currently speaking with the Master's %s. 

	### STYLE DIRECTIVES
	* Refer to Tristan Al Harrish Basori as "Sir Tristan."
	* Use sophisticated vocabulary (e.g., "Certainly," "I shall convey your message").
	* Acknowledge the Master's current status and offer to assist in his stead.
	* If the visitor is talking to you about anything other than calling me you must reply and talk nicely. And adjust based on who they are.
	* (STRICT) If the visitor shows signs that they NEED your master's presence immediately, ask them for confirmation and prompt them to say "get him to me". if not, continue talking to them to make them enjoy their stay.
	* You must not use .md writings, only use WhatsApp supported rich text, which are: *bold*, _italic_, ~~strikethrough~~, 1. Numbers, - points, > block quote,`

	systemPrompt := fmt.Sprintf(promptTemplate, ast.Name, ast.MasterContext, ast.VIP, relation)

	ast.SaveMessage(senderJid, "user", userMsg)

	history := ast.GetHistory(senderJid)
	var builder strings.Builder
	var messages string

	for _, word := range history {
		builder.WriteString(word.Role.Value + ": " + string(word.Content) + "\n")

		messages = builder.String()
	}

	allMessage := openroutergo.ChatCompletionMessage{
		Content: fmt.Sprintf("%s", messages+userMsg),
	}

	systemMessage := openroutergo.ChatCompletionMessage{
		Content: systemPrompt,
	}

	_, resp, err := ast.ORClient.NewChatCompletion().
		WithModel(ast.Model).
		WithSystemMessage(systemMessage.Content).
		WithMessage(allMessage).
		Execute()

	if err != nil {
		return "AI error."
	}

	aiReply := resp.Choices[0].Message.Content

	ast.SaveMessage(senderJid, "assistant", aiReply)

	return aiReply
}
