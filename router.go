package main

import (
	"fmt"

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
	adapter   *bluetooth.Adapter
}

func (c *BeaconCmd) Run() error {
	CustomUUID, _ = bluetooth.ParseUUID(CustomUUIDString)
	c.adapter = bluetooth.DefaultAdapter
	err := c.adapter.Enable()
	if err != nil {
		return fmt.Errorf("failed to enable adapter, is bluetooth on? error: %+v", err)
	}

	return c.Broadcast()

}

type ScanCmd struct {
	adapter bluetooth.Adapter
}

func (c *ScanCmd) Run() error {

	c.adapter = *bluetooth.DefaultAdapter
	err := c.adapter.Enable()
	if err != nil {
		return fmt.Errorf("failed to enable adapter, is bluetooth on? error: %+v", err)
	}

	return c.Scan()
}
