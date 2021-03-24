package main

import (
	"os"

	"github.com/cirss/blaze/pkg/blaze"
	"github.com/cirss/go-cli/pkg/cli"
)

var Program *cli.ProgramContext

func init() {
	Program = cli.NewProgramContext("blaze", main)
}

func main() {
	cc := blaze.NewBlazeCommandContext(Program)
	cc.InvokeCommand(os.Args)
}
