package cmd

import (
	"flag"
	"fmt"

	"github.com/tristnaja/clark/internal/whatsapp"
)

func ExecVIP(args []string, ast *whatsapp.Assistant) error {
	available, err := ast.CheckAst()

	if err != nil {
		return err
	}

	if !available {
		return fmt.Errorf("No assistant is initiated Sir. Do 'clark init' first.")
	}

	cmd := flag.NewFlagSet("vip", flag.ContinueOnError)
	var addTarget string
	var delTarget string

	cmd.StringVar(&addTarget, "add", "", "Add New VIP")
	cmd.StringVar(&addTarget, "a", "", "Add New VIP (shorthand)")
	cmd.StringVar(&delTarget, "delete", "", "Delete VIP")
	cmd.StringVar(&delTarget, "d", "", "Delete VIP (shorthand)")

	if err = cmd.Parse(args); err != nil {
		return fmt.Errorf("parsing args: %w", err)
	}

	if addTarget == "" && delTarget == "" {
		cmd.Usage()
		return fmt.Errorf("empty input")
	}

	if addTarget != "" {
		err = ast.VIP.AddVIP(addTarget)

		if err != nil {
			return fmt.Errorf("adding new VIP: %w", err)
		}

		fmt.Printf("\nAdded %v to our VIP list sir.\n", addTarget)
	}

	if delTarget != "" {
		err = ast.VIP.DeleteVIP(delTarget)

		if err != nil {
			return fmt.Errorf("adding new VIP: %w", err)
		}

		fmt.Printf("\nDeleted %v from our VIP list sir.\n", delTarget)
	}

	fmt.Println("\nNew VIP List:")
	for index, person := range ast.VIP.VIP {
		fmt.Printf("%v | %v\n", index, person)
	}

	return nil
}
