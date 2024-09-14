package macros

import (
	"fmt"

	"marlinraker/src/shared"
)

type sdcardPrintFileMacro struct{}

func (sdcardPrintFileMacro) Description() string {
	return "Loads a SD file and starts the print. May include files in subdirectories."
}

func (sdcardPrintFileMacro) Execute(manager *MacroManager, context shared.ExecutorContext, _ []string, _ Objects, params Params) error {
	fileName, err := params.RequireString("filename")
	if err != nil {
		return fmt.Errorf("filename: %w", err)
	}

	printManager := manager.printer.GetPrintManager()
	if err = printManager.SelectFile(fileName); err != nil {
		return fmt.Errorf("failed to select file %q: %w", fileName, err)
	}

	if err = printManager.Start(context); err != nil {
		return fmt.Errorf("failed to start print: %w", err)
	}
	return nil
}
