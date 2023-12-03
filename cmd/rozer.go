package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/alecthomas/participle/v2"
)

var (
	basicParser = participle.MustBuild[Program](
		participle.Lexer(pipelineLexer),
		participle.CaseInsensitive("Ident"),
		participle.Unquote("String"),
		participle.UseLookahead(4),
	)

	cli struct {
		File string `arg:"" type:"existingfile" help:"File to parse."`
	}
)

func Parse(r io.Reader) (*Program, error) {
	program, err := basicParser.Parse("", r)
	if err != nil {
		return nil, err
	}
	return program, nil
}

func main() {
	ctx := kong.Parse(&cli)
	r, err := os.Open(cli.File)
	ctx.FatalIfErrorf(err)
	defer r.Close()
	program, err := Parse(r)
	ctx.FatalIfErrorf(err)

	// repr.Println(program)
	fmt.Printf("%s", program.String())

	executable, errors := program.Compile()

	if errors.Len() > 0 {
		log.Printf("errors found during compilation:")
		for _, err := range *errors.errs {
			log.Printf("error: %s", err)
		}
		os.Exit(1)
	}

	executable.Execute(nil)
}
