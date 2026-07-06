package scan

import (
	"blueshare/internal/spec"
	"fmt"
	"sync"
	"time"

	"tinygo.org/x/bluetooth"
)

var msg = make([]byte, 256)
var mu = sync.Mutex{}

func StartScan() error {

	adapter := bluetooth.DefaultAdapter
	err := adapter.Enable()
	if err != nil {
		return fmt.Errorf("failed to enable adapter, is bluetooth on? error: %+v", err)
	}

	// mu := sync.Mutex

	go func(adapter *bluetooth.Adapter) {
		err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {

			// for _, data := range device.ServiceData() {
			// 	if data.UUID == spec.CustomUUID {
			// 		fmt.Println("data found: ", string(data.Data))
			// 	}
			// }

			for _, data := range device.ManufacturerData() {
				if data.CompanyID == 0xff {
					seqNum := int(data.Data[0])
					// sequence number holds the index of the packet being sent and each packet is of fixed size.
					baseIndex := seqNum * (spec.ByteSizeManufacturerData - 1)
					mu.Lock()
					copy(msg[baseIndex:], data.Data[1:])
					mu.Unlock()

				}
			}

		})

		fmt.Printf("error found while scanning, error : %+v", err)
	}(adapter)

	time.Sleep(time.Duration(spec.BroadcastWindow) * time.Second)

	defer adapter.StopScan()

	fmt.Println("the message is: ", string(msg))

	return nil

}
