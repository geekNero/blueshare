package main

import (
	"fmt"

	"tinygo.org/x/bluetooth"
)

func (c *BeaconCmd) Broadcast() error {

	advOptions := bluetooth.AdvertisementOptions{
		ServiceData: []bluetooth.ServiceDataElement{
			{
				UUID: CustomUUID,
				Data: []byte(c.Message),
			},
		},
	}

	adv := c.adapter.DefaultAdvertisement()
	err := adv.Configure(advOptions)
	if err != nil {
		return fmt.Errorf("failed to configure advertisement, error: %+v", err)
	}

	err = adv.Start()
	if err != nil {
		return fmt.Errorf("failed to start advertisement, error: %+v", err)
	}

	select {}
}
