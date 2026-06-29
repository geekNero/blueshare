package main

import (
	"github.com/alecthomas/kong"
	"tinygo.org/x/bluetooth"
)

var CustomUUID bluetooth.UUID

func main() {

	var cli CLI

	ctx := kong.Parse(&cli,
		kong.Name("blueshare"),
		kong.Description("A small text broadcast utility via bluetooth."),
	)

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
