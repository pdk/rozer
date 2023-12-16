package main

import (
	"github.com/alecthomas/participle/v2/lexer"
)

var (
	pipelineLexer = lexer.MustSimple([]lexer.SimpleRule{
		{
			Name:    "Comment",
			Pattern: `#[^\n]*`,
		},
		{
			Name:    "AndAnd",
			Pattern: `&&`,
		},
		{
			Name:    "OrOr",
			Pattern: `\|\|`,
		},
		{
			Name:    "PlusEqual",
			Pattern: `\+=`,
		},
		{
			Name:    "Plus",
			Pattern: `\+`,
		},
		{
			Name:    "Minus",
			Pattern: `-`,
		},
		{
			Name:    "Star",
			Pattern: `\*`,
		},
		{
			Name:    "Slash",
			Pattern: `/`,
		},
		{
			Name:    "Percent",
			Pattern: `%`,
		},
		{
			Name:    "EqualEqual",
			Pattern: `==`,
		},
		{
			Name:    "BangEqual",
			Pattern: `!=`,
		},
		{
			Name:    "Bang",
			Pattern: `!`,
		},
		{
			Name:    "MoreMoreMore",
			Pattern: `>>>`,
		},
		{
			Name:    "MoreMore",
			Pattern: `>>`,
		},
		{
			Name:    "LessEqual",
			Pattern: `<=`,
		},
		{
			Name:    "MoreEqual",
			Pattern: `>=`,
		},
		{
			Name:    "Less",
			Pattern: `<`,
		},
		{
			Name:    "More",
			Pattern: `>`,
		},
		{
			Name:    "Bang",
			Pattern: `!`,
		},
		{
			Name:    "ColonEqual",
			Pattern: `:=`,
		},
		{
			Name:    "Colon",
			Pattern: `:`,
		},
		{
			Name:    "String",
			Pattern: `"(\\"|[^"])*"`,
		},
		{
			Name:    "Function",
			Pattern: `fn`,
		},
		{
			Name:    "FTail",
			Pattern: `ƒ`,
		},
		{
			Name:    "Ident",
			Pattern: `[a-zA-Z_]\w*`,
		},
		{
			Name:    "Punct",
			Pattern: `[-[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`,
		},
		{
			Name:    "EOL",
			Pattern: `[\n\r]+`,
		},
		{
			Name:    "whitespace",
			Pattern: `[ \t]+`,
		},
		{
			// yyyy-mm-ddThh:mm:ss
			// yyyy-mm-ddThh:mm:ss.nnn
			// yyyy-mm-ddThh:mm:ss.nnn-10:00
			// yyyy-mm-ddThh:mm:ss+10:00
			Name:    "DateTime",
			Pattern: `\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\d(\.\d+)?([+-]\d\d:\d\d)?`,
		},
		{
			// yyyy-mm-dd
			Name:    "Date",
			Pattern: `\d\d\d\d-\d\d-\d\d`,
		},
		{
			// hh:mm:ss
			// hh:mm:ss.nnn
			Name:    "Time",
			Pattern: `\d\d:\d\d:\d\d(\.\d+)?`,
		},
		{
			// Y, y = years; M = months; D, d = days; H, h = hours; m = minutes; S, s = seconds
			// 2y => 2 years
			// 3M => 3 months
			// 3m => 3 minutes
			Name:    "TimeSpan",
			Pattern: `(\d+[ymdhs])+`,
		},
		{
			Name:    "Float",
			Pattern: `[-+]?\d+\.\d+`,
		},
		{
			Name:    "Integer",
			Pattern: `[-+]?\d\d*`,
		},
	})
)
