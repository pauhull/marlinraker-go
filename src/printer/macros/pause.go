package macros

import (
	"fmt"

	"marlinraker/src/shared"
)

type pauseMacro struct{}

func (pauseMacro) Description() string {
	return "Pauses the current print"
}

func (pauseMacro) Execute(manager *MacroManager, context shared.ExecutorContext, _ []string, _ Objects, _ Params) error {
	err := manager.printer.GetPrintManager().Pause(context)
	if err != nil {
		return fmt.Errorf("failed to pause print: %w", err)
	}
	return nil
}
