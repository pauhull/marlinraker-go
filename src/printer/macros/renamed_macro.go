package macros

import (
	"fmt"

	"marlinraker/src/shared"
)

type renamedMacro struct {
	description string
	original    Macro
}

func (macro renamedMacro) Description() string {
	return macro.description
}

func (macro renamedMacro) Execute(manager *MacroManager, context shared.ExecutorContext, rawParams []string, objects Objects, params Params) error {
	err := macro.original.Execute(manager, context, rawParams, objects, params)
	if err != nil {
		return fmt.Errorf("failed to execute macro: %w", err)
	}
	return nil
}
