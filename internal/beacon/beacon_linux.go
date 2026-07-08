//go:build linux

package beacon

import (
	"blueshare/internal/spec"
	"fmt"
	"sync"
	"time"

	"tinygo.org/x/bluetooth"
)

var (
	advertisementSlots   uint8
	advertisementStarted bool
	message              string
	mu                   sync.Mutex
)

func validateMessage(msg string) error {
	if len(msg) == 0 {
		return fmt.Errorf("message cannot be empty")
	}
	if len(msg) > spec.MaximumMessageLength {
		return fmt.Errorf("message length longer than %d chars", spec.MaximumMessageLength)
	}

	return nil
}

// windowedRoundRobin breaks the advertisement content into smaller packets and distributes them among
// the supported instances for advertisement in a round robin manner until the total time window expires.
func windowedRoundRobin(adapter *bluetooth.Adapter, c *BeaconCmd) {
	data := []byte(c.Message)
	payloadSize := spec.ByteSizeManufacturerData - 1
	numPackets := (len(c.Message) + payloadSize - 1) / payloadSize
	// timeBracket := (spec.BroadcastWindow / (numPackets * 2))
	timeBracket := spec.DefaultPerPacketDuration
	var sequenceNum uint8
	var index = 0
	for {
		wg := sync.WaitGroup{}
		for range advertisementSlots {
			sequenceNum = uint8(index % numPackets)
			offset := int(sequenceNum) * (payloadSize)
			endOffset := min(offset+payloadSize, len(c.Message))
			wg.Add(1)
			go func(adapter *bluetooth.Adapter, sequenceNum uint8, timeBracket int64, payload []byte, wg *sync.WaitGroup) {
				defer wg.Done()
				data := make([]byte, 1, spec.ByteSizeManufacturerData)
				data[0] = sequenceNum
				data = append(data, payload...)
				advertise(data, adapter, timeBracket)

			}(adapter, sequenceNum, int64(timeBracket)*int64(time.Second), data[offset:endOffset], &wg)
			index = (index + 1) % numPackets
		}
		wg.Wait()
	}
}

func advertise(data []byte, adapter *bluetooth.Adapter, duration int64) {

	fmt.Println("broadcasting new packet: ", string(data))
	defer fmt.Println("my packet was: ", string(data))

	advOptions := bluetooth.AdvertisementOptions{
		// LocalName: string(data),
		ManufacturerData: []bluetooth.ManufacturerDataElement{
			{
				CompanyID: 0xff,
				Data:      data,
			},
		},
		// ServiceData: []bluetooth.ServiceDataElement{
		// 	{
		// 		UUID: spec.CustomUUID,
		// 		Data: []byte(message.msg),
		// 	},
		// },
	}
	mu.Lock()
	advertisement := adapter.NewAdvertisement()
	err := advertisement.Configure(advOptions)
	if err != nil {
		mu.Unlock()
		fmt.Printf("failed to configure advertisement, error: %+v", err)
		return
	}

	mu.Lock()
	err = advertisement.Start()
	if err != nil {
		mu.Unlock()
		fmt.Printf("failed to start advertisement, error: %+v", err)
		return
	}
	mu.Unlock()

	time.Sleep(time.Duration(duration))
	mu.Lock()
	advertisement.Stop()
	mu.Unlock()
}

func Broadcast(c *BeaconCmd) spec.ErrorResponse {

	err := validateMessage(c.Message)
	if err != nil {
		return spec.ErrorResponse{
			Error: spec.ErrInvalidMessage,
			Err:   err,
		}
	}

	adapter := bluetooth.DefaultAdapter
	err = adapter.Enable()
	if err != nil {
		return spec.ErrorResponse{
			Error: spec.ErrFailedToEnableAdapter,
			Err:   err,
		}
	}

	advertisementSlots, err = adapter.AvailableAdvertisementSlots()
	if err != nil {
		return spec.ErrorResponse{
			Error: spec.ErrUnknown,
			Err:   err,
		}
	}

	myAddress, _ := adapter.Address()
	fmt.Println("broadcasting from address: ", myAddress.String())

	go windowedRoundRobin(adapter, c)

	return spec.ErrorResponse{}
}
