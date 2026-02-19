package cmd

import (
	"github.com/tristnaja/clark/internal/whatsapp"
)

func ExecInit(ast *whatsapp.Assistant) error {
	err := ast.AstSettingInit()

	if err != nil {
		return err
	}

	return nil
}
