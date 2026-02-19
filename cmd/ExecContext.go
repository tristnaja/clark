package cmd

import (
	"flag"
	"fmt"

	"github.com/tristnaja/clark/internal/whatsapp"
)

func ExecContext(args []string, ast *whatsapp.Assistant) error {
	available, err := ast.CheckAst()

	if err != nil {
		return err
	}

	if !available {
		return fmt.Errorf("No assistant is initiated Sir. Do 'clark init' first.")
	}

	cmd := flag.NewFlagSet("ctx", flag.ContinueOnError)
	var contex string

	cmd.StringVar(&contex, "change", "", "Change Context")
	cmd.StringVar(&contex, "c", "", "Change Context (Shorthand)")

	err = cmd.Parse(args)

	if err != nil {
		return err
	}

	err = ast.SetMasterContext(contex)

	if err != nil {
		return err
	}

	return nil
}
