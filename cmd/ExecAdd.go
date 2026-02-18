package cmd

import (
	"flag"
	"fmt"

	"github.com/tristnaja/clark/internal"
)

func ExecAdd(args []string, ast *internal.Assistant) error {
	cmd := flag.NewFlagSet("add", flag.ContinueOnError)
	var newVIP string

	cmd.StringVar(&newVIP, "vip", "", "New VIP")
	cmd.StringVar(&newVIP, "v", "", "New VIP (shorthand)")

	err := cmd.Parse(args)

	if err != nil {
		return fmt.Errorf("parsing args: %w", err)
	}

	available, err := ast.CheckAst()

	if err != nil {
		return err
	}

	if !available {
		return fmt.Errorf("No assistant is initiated Sir. Do 'clark init' first.")
	}

	if newVIP == "" {
		cmd.Usage()
		return fmt.Errorf("empty input")
	}

	err = ast.VIP.AddVIP(newVIP)

	if err != nil {
		return fmt.Errorf("adding new VIP: %w", err)
	}

	fmt.Printf("\nAdded %v to our VIP list sir.\n", newVIP)

	fmt.Println("\nNew VIP List:")
	for index, person := range ast.VIP.VIP {
		fmt.Printf("%v | %v\n", index, person)
	}

	return nil
}
