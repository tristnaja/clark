package cmd

import (
	"fmt"

	"github.com/tristnaja/clark/internal"
)

func ExecView(ast *internal.Assistant) {
	fmt.Println("Here is your settings, Sir:")
	fmt.Println("Assistant Name:", ast.Name)
	fmt.Println("Assistant Model:", ast.Model)
	fmt.Println("Active Status:", ast.Status)
	fmt.Println("Master Context:", ast.MasterContext)
	fmt.Println("\nHere is your VIP list:")
	for jid, name := range ast.VIP.VIP {
		fmt.Printf("%v | %v\n", jid, name)
	}
}
