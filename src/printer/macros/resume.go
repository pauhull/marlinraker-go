package macros

import (
	"fmt"

	"marlinraker/src/shared"
)

type resumeMacro struct{}

func (resumeMacro) Description() string {
	return "Resumes the print from a pause"
}

func (resumeMacro) Execute(manager *MacroManager, context shared.ExecutorContext, _ []string, _ Objects, _ Params) error {
	err := manager.printer.GetPrintManager().Resume(context)
	if err != nil {
		return fmt.Errorf("failed to resume print: %w", err)
	}
	return nil
}
