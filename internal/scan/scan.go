package scan

import (
	"blueshare/internal/spec"
	"fmt"
	"sync"

	"tinygo.org/x/bluetooth"
)

var (
	message  = make([]byte, 256)
	mu       = sync.Mutex{}
	messages = sync.Map{}
	adapter  *bluetooth.Adapter
)

func StartScan() spec.ErrorResponse {

	if adapter != nil {
		return spec.ErrorResponse{
			Error: spec.ErrScanAlreadyInProgress,
		}
	}

	adapter = bluetooth.DefaultAdapter
	err := adapter.Enable()
	if err != nil {
		return spec.ErrorResponse{
			Error:  spec.ErrFailedToEnableAdapter,
			ErrMsg: err.Error(),
		}
	}

	go func(adapter *bluetooth.Adapter) {
		err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {

			// for _, data := range device.ServiceData() {
			// 	if data.UUID == spec.CustomUUID {
			// 		fmt.Println("data found: ", string(data.Data))
			// 	}
			// }

			for _, data := range device.ManufacturerData() {
				if data.CompanyID == 0xff {
					// msg, _ := messages[device.A]
					seqNum := int(data.Data[0])
					// sequence number holds the index of the packet being sent and each packet is of fixed size.
					baseIndex := seqNum * (spec.ByteSizeManufacturerData - 1)
					for {
						emptyByteArray := make([]byte, spec.MaximumMessageLength)
						raw, _ := messages.LoadOrStore(device.Address.String(), &emptyByteArray)
						var message = raw.(*[]byte)

						copy((*message)[baseIndex:], data.Data[1:])
						ok := messages.CompareAndSwap(device.Address.String(), raw, message)
						if ok {
							break
						}
					}

				}
			}

		})

		if err != nil {
			fmt.Printf("error found while scanning, error : %+v", err)
		}
	}(adapter)

	return spec.ErrorResponse{}
}

func StopScan() error {
	if adapter == nil {
		return nil
	}

	err := adapter.StopScan()
	if err != nil {
		return fmt.Errorf("failed to stop scanning, error: %+v", err)
	}

	adapter = nil
	return nil
}

func FetchMessages() *[]string {
	replies := make([]string, 0)

	messages.Range(func(key, value any) bool {
		reply := value.(*[]byte)
		replies = append(replies, string(*reply))
		return true
	})
	return &replies
}
