//go:build linux

package beacon

import (
	"blueshare/internal/spec"
	"fmt"

	"tinygo.org/x/bluetooth"
)

func Broadcast(message *Message) error {

	adapter := bluetooth.DefaultAdapter
	err := adapter.Enable()
	if err != nil {
		return fmt.Errorf("failed to enable adapter, is bluetooth on? error: %+v", err)
	}

	myAddress, _ := adapter.Address()

	fmt.Println("broadcasting from address: ", myAddress.String())

	advOptions := bluetooth.AdvertisementOptions{
		ServiceData: []bluetooth.ServiceDataElement{
			{
				UUID: spec.CustomUUID,
				Data: []byte(message.msg),
			},
		},
	}

	adv := adapter.DefaultAdvertisement()
	err = adv.Configure(advOptions)
	if err != nil {
		return fmt.Errorf("failed to configure advertisement, error: %+v", err)
	}

	err = adv.Start()
	if err != nil {
		return fmt.Errorf("failed to start advertisement, error: %+v", err)
	}

	select {}
}
