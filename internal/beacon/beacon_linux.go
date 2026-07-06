//go:build linux

package beacon

import (
	"blueshare/internal/spec"
	"fmt"
	"sync"
	"time"

	"tinygo.org/x/bluetooth"
)

var advertisementSlots uint8

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
func windowedRoundRobin(adapter *bluetooth.Adapter, message *Message) error {
	data := []byte(message.msg)
	payloadSize := spec.ByteSizeManufacturerData - 1
	numPackets := (len(message.msg) + payloadSize - 1) / payloadSize
	timeBracket := (spec.BroadcastWindow / (numPackets * 2))
	var sequenceNum uint8
	var index = 0
	for index < (numPackets * 2) {
		wg := sync.WaitGroup{}
		for range advertisementSlots {
			if index >= (numPackets * 2) {
				break
			}
			sequenceNum = uint8(index % numPackets)
			offset := int(sequenceNum) * (payloadSize)
			endOffset := min(offset+payloadSize, len(message.msg))
			wg.Add(1)
			go func(adapter *bluetooth.Adapter, sequenceNum uint8, timeBracket int64, payload []byte, wg *sync.WaitGroup) {
				defer wg.Done()
				data := make([]byte, 1, spec.ByteSizeManufacturerData)
				data[0] = sequenceNum
				data = append(data, payload...)
				advertise(data, adapter, timeBracket)

			}(adapter, sequenceNum, int64(timeBracket)*int64(time.Second), data[offset:endOffset], &wg)
			index++
		}
		wg.Wait()
	}

	return nil
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
	adv := adapter.NewAdvertisement()
	err := adv.Configure(advOptions)
	if err != nil {
		fmt.Printf("failed to configure advertisement, error: %+v", err)
		return
	}

	err = adv.Start()
	if err != nil {
		fmt.Printf("failed to start advertisement, error: %+v", err)
		return
	}
	defer adv.Stop()

	time.Sleep(time.Duration(duration))

}

func Broadcast(message *Message) error {

	err := validateMessage(message.msg)
	if err != nil {
		return fmt.Errorf("message format incorrect, error: %s", err.Error())
	}

	adapter := bluetooth.DefaultAdapter
	err = adapter.Enable()
	if err != nil {
		return fmt.Errorf("failed to enable adapter, is bluetooth on? error: %+v", err)
	}

	advertisementSlots, err = adapter.GetAdvertisementSlots()
	if err != nil {
		return fmt.Errorf("failed to get available advertisement slots, error: %+v", err)
	}

	myAddress, _ := adapter.Address()
	fmt.Println("broadcasting from address: ", myAddress.String())

	err = windowedRoundRobin(adapter, message)
	if err != nil {
		return fmt.Errorf("error while using windowedRoundRobin, error: %+v", err)
	}

	return nil
}
