package scan

import (
	"blueshare/internal/spec"
	"fmt"

	"tinygo.org/x/bluetooth"
)

// var msgs = []string{}

func Scan() error {

	adapter := bluetooth.DefaultAdapter
	err := adapter.Enable()
	if err != nil {
		return fmt.Errorf("failed to enable adapter, is bluetooth on? error: %+v", err)
	}

	err = adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {

		for _, data := range device.ServiceData() {
			if data.UUID == spec.CustomUUID {
				fmt.Println("data found: ", string(data.Data))
			}
		}

	})

	defer adapter.StopScan()
	if err != nil {
		return fmt.Errorf("failed to start scanning, error: %+v", err)
	}

	return nil

}
