package main

import (
	"fmt"

	"tinygo.org/x/bluetooth"
)

// var msgs = []string{}

func (c *ScanCmd) Scan() error {

	err := c.adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {

		if !device.HasServiceUUID(CustomUUID) {
			return
		}

		for _, data := range device.ServiceData() {
			if data.UUID == CustomUUID {
				fmt.Println("data found: ", string(data.Data))
			}
		}

	})

	defer c.adapter.StopScan()
	if err != nil {
		return fmt.Errorf("failed to start scanning, error: %+v", err)
	}

	return nil

}
