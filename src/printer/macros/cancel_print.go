package macros

import (
	"fmt"

	"marlinraker/src/shared"
)

type cancelPrintMacro struct{}

func (cancelPrintMacro) Description() string {
	return "Cancel the current print"
}

func (cancelPrintMacro) Execute(manager *MacroManager, context shared.ExecutorContext, _ []string, _ Objects, _ Params) error {
	err := manager.printer.GetPrintManager().Cancel(context)
	if err != nil {
		return fmt.Errorf("failed to cancel print: %w", err)
	}
	return nil
}
