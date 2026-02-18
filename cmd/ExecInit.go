package cmd

import "github.com/tristnaja/clark/internal"

func ExecInit(ast *internal.Assistant) error {
	err := ast.AstSettingInit()

	if err != nil {
		return err
	}

	return nil
}
