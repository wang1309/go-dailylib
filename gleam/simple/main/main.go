package main

import (
	"dailylib/gleam/simple"
	"dailylib/gleam/simple/file"
	"dailylib/gleam/simple/flow"
	"flag"
	"strings"
)

var (
	Tokenize      = simple.RegisterMapper(tokenize)
	AppendOne     = simple.RegisterMapper(appendOne)
	Sum           = simple.RegisterReducer(sum)
)

func tokenize(row []interface{}) error {
	line := simple.ToString(row[0])
	for _, s := range strings.FieldsFunc(line, func(r rune) bool {
		return !('A' <= r && r <= 'Z' || 'a' <= r && r <= 'z' || '0' <= r && r <= '9')
	}) {
		simple.Emit(s)
	}
	return nil
}

func appendOne(row []interface{}) error {
	row = append(row, 1)
	simple.Emit(row...)
	return nil
}

func sum(x, y interface{}) (interface{}, error) {
	return simple.ToInt64(x) + simple.ToInt64(y), nil
}

func main() {
	simple.Init() // If the command line invokes the mapper or reducer, execute it and exit.
	flag.Parse()  // optional, since gio.Init() will call this also.

	f := flow.New("top5 words in passwd").
		Read(file.Txt("/etc/hosts", 2)). // read a txt file and partitioned to 2 shards
		Map("tokenize", Tokenize). // invoke the registered "tokenize" mapper function.
		Map("appendOne", AppendOne). // invoke the registered "appendOne" mapper function.
		ReduceByKey("sum", Sum). // invoke the registered "sum" reducer function.
		Sort("sortBySum", flow.OrderBy(2, true)).
		Top("top5", 5, flow.OrderBy(2, false)).
		Printlnf("%s\t%d")

	f.Run()
}
