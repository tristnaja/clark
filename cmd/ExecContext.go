package cmd

import (
	"flag"
	"fmt"

	"github.com/tristnaja/clark/internal"
)

func ExecContext(args []string, ast *internal.Assistant) error {
	cmd := flag.NewFlagSet("ctx", flag.ContinueOnError)
	var contex string

	cmd.StringVar(&contex, "ctx", "", "Add Context")
	cmd.StringVar(&contex, "c", "", "Add Context (Shorthand)")

	err := cmd.Parse(args)

	if err != nil {
		return err
	}

	available, err := ast.CheckAst()

	if err != nil {
		return err
	}

	if !available {
		return fmt.Errorf("No assistant is initiated Sir. Do 'clark init' first.")
	}

	err = ast.SetMasterContext(contex)

	if err != nil {
		return err
	}

	return nil
}
