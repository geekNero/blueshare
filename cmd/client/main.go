package main

import (
	"blueshare/internal/beacon"
	"blueshare/internal/scan"

	"github.com/alecthomas/kong"
)

type CLI struct {
	Buzz beacon.BeaconCmd `cmd:"" help:"Broadcast a message."`
	Scan scan.ScanCmd     `cmd:"" help:"Scan for broadcasted messages."`
}

func main() {

	var cli CLI

	ctx := kong.Parse(&cli,
		kong.Name("blueshare"),
		kong.Description("A small text broadcast utility via bluetooth."),
	)

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
