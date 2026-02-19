package main

import (
	"errors"
	"log"
	"os"

	"github.com/tristnaja/clark/cmd"
	"github.com/tristnaja/clark/internal/whatsapp"
)

func main() {
	log.SetPrefix("clark: ")
	log.SetFlags(0)
	var err error
	commands := map[string]struct{}{
		"init":   {},
		"run":    {},
		"vip":    {},
		"ctx":    {},
		"toggle": {},
		"view":   {},
	}

	if len(os.Args) < 2 {
		log.Fatal("usage: clark [cmd]")
	}

	if _, exist := commands[os.Args[1]]; !exist {
		log.Fatalf("unknown command '%v'", os.Args[1])
	}

	if len(os.Args) > 2 && os.Args[1] == "run" {
		log.Fatal("unnecessary argument(s), usage: clark run")
	}

	if len(os.Args) < 3 && (os.Args[1] == "add" || os.Args[1] == "ctx") {
		log.Fatalf("usage: clark %v [args]", os.Args[1])
	}

	ast, err := whatsapp.AssistantInit()

	if err != nil {
		log.Fatalf("fail to create assistant: %v", err)
	}

	switch os.Args[1] {
	case "init":
		err = cmd.ExecInit(ast)
	case "run":
		err = cmd.ExecRun(ast)
	case "vip":
		err = cmd.ExecVIP(os.Args[2:], ast)
	case "ctx":
		err = cmd.ExecContext(os.Args[2:], ast)
	case "toggle":
		err = cmd.ExecToggle(ast)
	case "view":
		cmd.ExecView(ast)
	default:
		err = errors.New("unknown command sir, here are the commands: init, run, add, ctx, toggle")
	}

	if err != nil {
		log.Fatal(err)
	}

	ast.DB.DB.Close()
}
