package cmd

import (
	"fmt"

	"github.com/tristnaja/clark/internal"
)

func ExecToggle(ast *internal.Assistant) error {
	available, err := ast.CheckAst()

	if err != nil {
		return err
	}

	if !available {
		return fmt.Errorf("No assistant is initiated Sir. Do 'clark init' first.")
	}

	err = ast.ToggleStatus()

	if err != nil {
		return err
	}

	return nil
}
