package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	"github.com/tristnaja/clark/internal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func ExecRun(ast *internal.Assistant) error {
	available, err := ast.CheckAst()

	if err != nil {
		return err
	}

	if !available {
		return fmt.Errorf("No assistant is initiated Sir. Do 'clark init' first.")
	}

	if ast.MasterContext == "" {
		return fmt.Errorf("No context yet Sir. Do 'clark ctx -c [context]' first.")
	}

	if ast.Status == false {
		return fmt.Errorf("Clark is not active yet Sir. Do 'clark toggle' to toggle it on.")
	}

	fmt.Println("| Assistant Active")
	fmt.Printf("Assistant Name: %v\n", ast.Name)
	fmt.Printf("Master Context: %v\n", ast.MasterContext)

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New(context.Background(), "sqlite3", "file:mystore.db?_foreign_keys=on", dbLog)

	if err != nil {
		return fmt.Errorf("fail to initiate database container: %v", err)
	}

	defer container.Close()

	rawDb, err := sql.Open("sqlite3", "mystore.db")

	if err != nil {
		return fmt.Errorf("fail to open database: %v", err)
	}

	defer rawDb.Close()

	if err := rawDb.Ping(); err != nil {
		return fmt.Errorf("fail to ping database: %v", err)
	}

	deviceStore, err := container.GetFirstDevice(context.Background())

	if err != nil {
		return fmt.Errorf("fail to get device connection: %v", err)
	}

	client := whatsmeow.NewClient(deviceStore, waLog.Stdout("Client", "INFO", true))
	client.AddEventHandler(internal.EventHandler(client, ast))

	if client.Store.ID == nil {
		fmt.Println("No session found. Please scan QR code.")
		qrChan, _ := client.GetQRChannel(context.Background())
		err := client.Connect()

		if err != nil {
			return fmt.Errorf("fail to connect through QR: %v", err)
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			}
		}
	} else {
		fmt.Println("Existing session found. Connecting...")
		err := client.Connect()

		if err != nil {
			return fmt.Errorf("fail to connect to existing session: %v", err)
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	client.Disconnect()
	return nil
}
