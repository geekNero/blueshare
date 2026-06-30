package main

import (
	"blueshare/internal/beacon"
	"blueshare/internal/scan"
	"blueshare/internal/spec"

	"tinygo.org/x/bluetooth"
)

type CLI struct {
	Buzz BeaconCmd `cmd:"" help:"Broadcast a message."`
	Scan ScanCmd   `cmd:"" help:"Scan for broadcasted messages."`
}

type BeaconCmd struct {
	Message   string `arg:"" name:"msg" help:"String to be broadcast."`
	Frequency int    `default:"2" name:"frequency" short:"f" help:"Frequency at which the message should be broardcast.\nWithin range 1-3"`
	Once      bool   `help:"Only broadcast message once."`
}

func (c *BeaconCmd) Run() error {
	spec.CustomUUID, _ = bluetooth.ParseUUID(spec.CustomUUIDString)
	return beacon.Broadcast(beacon.NewMessage(c.Message, c.Frequency))

}

type ScanCmd struct {
}

func (c *ScanCmd) Run() error {

	return scan.Scan()
}
