package macros

import (
	"fmt"

	"marlinraker/src/shared"
)

type restoreGcodeState struct {
}

func (restoreGcodeState) Description() string {
	return "Restore a previously saved G-Code state"
}

func (restoreGcodeState) Execute(manager *MacroManager, context shared.ExecutorContext, _ []string, _ Objects, params Params) error {
	name, exists := params["NAME"]
	if !exists {
		name = "default"
	}
	err := manager.printer.RestoreGcodeState(context, name)
	if err != nil {
		return fmt.Errorf("failed to restore G-code state: %w", err)
	}
	return nil
}
