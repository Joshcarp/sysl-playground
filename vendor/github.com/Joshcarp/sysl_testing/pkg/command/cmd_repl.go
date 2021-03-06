package command

import (
	"os"

	"github.com/Joshcarp/sysl_testing/pkg/eval"
	"gopkg.in/alecthomas/kingpin.v2"
)

type replCmd struct{}

func (p *replCmd) Name() string            { return "repl" }
func (p *replCmd) RequireSyslModule() bool { return false }

func (p *replCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Enter a sysl REPL")
	return cmd
}

func (p *replCmd) Execute(args ExecuteArgs) error {
	s := &eval.Scope{}
	repl := eval.NewREPL(os.Stdin, os.Stdout)
	for {
		if err := repl(s, nil, nil); err != nil {
			return nil // means EOF
		}
	}
}
