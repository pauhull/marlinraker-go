package macros

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	log "github.com/sirupsen/logrus"

	"marlinraker/src/shared"
)

type customMacro struct {
	description string
	tmpl        *template.Template
}

func newCustomMacro(name string, description string, content string) (*customMacro, error) {
	tmpl, err := template.New(name).Funcs(sprig.FuncMap()).Parse(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}
	return &customMacro{description: description, tmpl: tmpl}, nil
}

func (macro *customMacro) Description() string {
	return macro.description
}

type macroContext struct {
	manager   *MacroManager
	Printer   Objects
	Params    Params
	RawParams []string
}

func (context macroContext) ActionRespondInfo(message string) string {
	if err := context.manager.printer.Respond("// " + message); err != nil {
		log.Errorf("Failed to respond: %v", err)
	}
	return ""
}

func (context macroContext) ActionRaiseError(message string) string {
	if err := context.manager.printer.Respond("!! " + message); err != nil {
		log.Errorf("Failed to respond: %v", err)
	}
	return ""
}

func (macro *customMacro) Execute(manager *MacroManager, context shared.ExecutorContext, rawParams []string, objects Objects, params Params) error {

	ctx := macroContext{
		manager:   manager,
		Printer:   objects,
		Params:    params,
		RawParams: rawParams,
	}

	builder := strings.Builder{}
	if err := macro.tmpl.Execute(&builder, ctx); err != nil {
		return fmt.Errorf("failed to execute macro: %w", err)
	}
	gcode := builder.String()
	<-context.QueueGcode(gcode, true)
	return nil
}
