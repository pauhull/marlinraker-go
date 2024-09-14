package macros

import (
	"fmt"

	"marlinraker/src/shared"
)

type sdcardResetFileMacro struct{}

func (sdcardResetFileMacro) Description() string {
	return "Clears a loaded SD File. Stops the print if necessary."
}

func (sdcardResetFileMacro) Execute(manager *MacroManager, context shared.ExecutorContext, _ []string, _ Objects, _ Params) error {
	err := manager.printer.GetPrintManager().Reset(context)
	if err != nil {
		return fmt.Errorf("failed to reset SD print: %w", err)
	}
	return nil
}
