package scan

import (
	"fmt"

	"tinygo.org/x/bluetooth"
)

var msg string

func Scan() error {

	adapter := bluetooth.DefaultAdapter
	err := adapter.Enable()
	if err != nil {
		return fmt.Errorf("failed to enable adapter, is bluetooth on? error: %+v", err)
	}

	ch := make(chan bool)

	// mu := sync.Mutex

	go func(ch chan bool, adapter *bluetooth.Adapter) {
		err = adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {

			// for _, data := range device.ServiceData() {
			// 	if data.UUID == spec.CustomUUID {
			// 		fmt.Println("data found: ", string(data.Data))
			// 	}
			// }

			for _, data := range device.ManufacturerData() {
				if data.CompanyID == 0xff {

					if data.Data[0] == 0x00 {
						ch <- true
						return
					}

					msg += string(data.Data)
				}
			}

		})

	}(ch, adapter)

	<-ch

	defer adapter.StopScan()
	if err != nil {
		return fmt.Errorf("failed to start scanning, error: %+v", err)
	}

	fmt.Println("the message is: ", msg)

	return nil

}
