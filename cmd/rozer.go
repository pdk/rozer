package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/pdk/rozer/lang"
)

var (
	basicParser = participle.MustBuild[lang.Program](
		participle.Lexer(lang.PipelineLexer),
		participle.CaseInsensitive("Ident"),
		participle.Unquote("String"),
		participle.UseLookahead(4),
	)

	cli struct {
		File string `arg:"" type:"existingfile" help:"File to parse."`
	}
)

func Parse(r io.Reader) (*lang.Program, error) {
	program, err := basicParser.Parse("", r)
	if err != nil {
		return nil, err
	}
	return program, nil
}

func main() {
	// ctx := kong.Parse(&cli)
	// r, err := os.Open(cli.File)
	// ctx.FatalIfErrorf(err)
	// defer r.Close()

	r, err := os.Open("/Users/pdk/src/rozer/short.roz")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	program, err := Parse(r)
	if err != nil {
		log.Fatal(err)
	}

	// repr.Println(program)
	fmt.Printf("%s", program.String())

	globalTypeMap := lang.TypeMap{}

	log.Printf("compiling...")
	executableProgram, errors := program.Compile(globalTypeMap)

	if errors.Len() > 0 {
		log.Printf("errors found during compilation:")
		for _, err := range *errors.Errs {
			log.Printf("error: %s", err)
		}
		os.Exit(1)
	}

	// log.Printf("here are the functions: ")
	// executableProgram.DumpFunctions()

	// log.Printf("here is the program: ")
	// executableProgram.DumpProgram()

	log.Printf("executing...")

	executableProgram.ExecuteProgram()
}
