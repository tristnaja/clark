package internal

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/eduardolat/openroutergo"
	"github.com/gen2brain/beeep"
	"github.com/joho/godotenv"
)

type Assistant struct {
	Name          string
	Status        bool
	MasterContext string
	VIP           *VIP
	DB            *Database
	ORClient      *openroutergo.Client
	Model         string
}

//go:embed utils/newPrompt.md
var promptTemplate string

func AssistantInit() (*Assistant, error) {
	err := godotenv.Load()

	if err != nil {
		return nil, fmt.Errorf("fail to load .env: %v", err)
	}

	apiKey := os.Getenv("OPENROUTER_API")
	client, err := openroutergo.NewClient().WithAPIKey(apiKey).Create()
	beeep.AppName = "Clark"

	if err != nil {
		return nil, fmt.Errorf("fail to initiate AI client: %w", err)
	}

	db, err := InitDB()

	if err != nil {
		return nil, err
	}

	vip := InitVIP(db)

	ast := &Assistant{
		Name:          "",
		Status:        false,
		MasterContext: "",
		VIP:           vip,
		DB:            db,
		ORClient:      client,
		Model:         "stepfun/step-3.5-flash:free",
	}

	err = ast.VIP.LoadVIP()

	if err != nil {
		return nil, err
	}

	ast.loadAssistant()
	return ast, nil
}

func (ast *Assistant) ToggleStatus() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var newStatus string

	ast.DB.DB.QueryRowContext(ctx, "SELECT value FROM assistant_setting WHERE key = 'status'").Scan(&newStatus)

	statusBool, err := strconv.ParseBool(newStatus)

	if err != nil {
		return err

	}
	newStatusStr := fmt.Sprintf("%v", !statusBool)

	query := `INSERT OR REPLACE INTO assistant_setting (key, value) VALUES (?, ?)`
	_, err = ast.DB.DB.ExecContext(ctx, query, "status", newStatusStr)

	if err != nil {
		return fmt.Errorf("failed to update status in DB: %w", err)
	}

	ast.Status = !statusBool

	fmt.Println("Current Status:")
	fmt.Println(ast.Status)

	return nil
}

func (ast *Assistant) SetMasterContext(contextInput string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	query := `INSERT OR REPLACE INTO assistant_setting (key, value) VALUES (?, ?)`
	_, err := ast.DB.DB.ExecContext(ctx, query, "context", contextInput)

	if err != nil {
		return err
	}

	ast.MasterContext = contextInput

	fmt.Println("Current Master Context:")
	fmt.Println(ast.MasterContext)

	return nil
}

func (ast *Assistant) AstSettingInit() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	query := `INSERT OR IGNORE INTO assistant_setting (key, value) VALUES (?, ?)`

	defaults := map[string]string{
		"name":    "Clark",
		"status":  "false",
		"context": "",
	}

	for key, value := range defaults {
		_, err := ast.DB.DB.ExecContext(ctx, query, key, value)

		if err != nil {
			return fmt.Errorf("fail to initialize default for %s: %w", key, err)
		}
	}

	return nil
}

func (ast *Assistant) CheckAst() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var count int
	query := `SELECT COUNT(*) FROM assistant_setting`

	err := ast.DB.DB.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("fail to load table <vip>: %w", err)
	}

	return count == 3, nil
}

func (ast *Assistant) loadAssistant() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var name, status, context string
	ast.DB.DB.QueryRowContext(ctx, "SELECT value FROM assistant_setting WHERE key = 'name'").Scan(&name)
	ast.DB.DB.QueryRowContext(ctx, "SELECT value FROM assistant_setting WHERE key = 'status'").Scan(&status)
	ast.DB.DB.QueryRowContext(ctx, "SELECT value FROM assistant_setting WHERE key = 'context'").Scan(&context)

	if name == "" {
		return fmt.Errorf("Name is empty. Error occured Sir.")
	}

	if context == "" {
		return fmt.Errorf("Master Context is empty. Error occured Sir.")
	}

	statusBool, err := strconv.ParseBool(status)
	if err != nil {
		return fmt.Errorf("Invalid status value Sir. Error: %w", err)
	}

	ast.Status = statusBool
	ast.Name = name
	ast.MasterContext = context

	return nil
}

func (ast *Assistant) GetHistory(jid string) ([]openroutergo.ChatCompletionMessage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	getHistoryQuery := "SELECT role, content FROM chat_history WHERE jid = ? ORDER BY timestamp ASC"
	rows, err := ast.DB.DB.QueryContext(ctx, getHistoryQuery, jid)

	if err != nil {
		return nil, fmt.Errorf("failed to query history: %w", err)
	}

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
	return history, nil
}

func (ast *Assistant) GetAIResponse(senderJid, userMsg string) (string, error) {
	if senderJid == "" {
		return "", fmt.Errorf("empty sender JID")
	}

	relation, isVIP := ast.VIP.CheckVIP(senderJid)

	if !isVIP {
		return "", fmt.Errorf("sender not in VIP list")
	}

	if userMsg == "" {
		return "", fmt.Errorf("empty message content")
	}

	err := ast.DB.SaveMessage(senderJid, "user", userMsg)

	if err != nil {
		return "", err
	}

	history, err := ast.GetHistory(senderJid)

	if err != nil {
		return "", err
	}

	if len(history) == 0 {
		return "", fmt.Errorf("no chat history available")
	}

	systemPrompt := fmt.Sprintf(promptTemplate, ast.Name, ast.MasterContext, ast.VIP, relation)

	query := ast.ORClient.NewChatCompletion().
		WithModel(ast.Model).
		WithSystemMessage(systemPrompt)

	for _, msg := range history {
		content := string(msg.Content)
		switch msg.Role.Value {
		case "user":
			query.WithUserMessage(content)
		case "assistant":
			query.WithAssistantMessage(content)
		}
	}

	_, resp, err := query.Execute()
	if err != nil {
		return "", fmt.Errorf("failed to execute model: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from model")
	}

	aiReply := resp.Choices[0].Message.Content

	err = ast.DB.SaveMessage(senderJid, "assistant", aiReply)

	if err != nil {
		return "", err
	}

	return aiReply, nil
}
