//go:build linux

package beacon

import (
	"blueshare/internal/spec"
	"fmt"
	"time"

	"tinygo.org/x/bluetooth"
)

func validateMessage(msg string) error {
	if len(msg) > spec.MaximumMessageLength {
		return fmt.Errorf("message length longer than %d chars", spec.MaximumMessageLength)
	}

	return nil
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

	myAddress, _ := adapter.Address()

	fmt.Println("broadcasting from address: ", myAddress.String())

	advertise([]byte(message.msg)[:5], adapter)
	advertise([]byte(message.msg)[6:], adapter)

	advertise([]byte{0x00}, adapter)

	// select {}
	return nil
}

func advertise(data []byte, adapter *bluetooth.Adapter) {

	advOptions := bluetooth.AdvertisementOptions{
		LocalName: "Blueshare",
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
	adv := adapter.DefaultAdvertisement()
	err := adv.Configure(advOptions)
	if err != nil {
		fmt.Printf("failed to configure advertisement, error: %+v", err)
	}

	err = adv.Start()
	if err != nil {
		fmt.Printf("failed to start advertisement, error: %+v", err)
	}

	time.Sleep(3 * time.Second)

	adv.Stop()

}
