//go:build linux

package beacon

import (
	"blueshare/internal/spec"
	"context"
	"fmt"
	"sync"
	"time"

	"tinygo.org/x/bluetooth"
)

type Window struct {
	advPool    []*bluetooth.Advertisement
	startIndex int
	endIndex   int
	set        bool
	cancel     context.CancelFunc
	data       []byte
	mu         sync.Mutex
	await      chan struct{}
	err        error
}

type Adapter struct {
	adapter *bluetooth.Adapter
	mu      sync.Mutex
}

var (
	window             Window
	advertisementSlots uint8
	adapter            Adapter
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

func setAdapter() error {
	adapter.mu.Lock()
	defer adapter.mu.Unlock()
	if adapter.adapter != nil {
		return nil
	}
	adapter.adapter = bluetooth.DefaultAdapter

	return adapter.adapter.Enable()
}

func setWindowLocked(msg string) {
	window.data = []byte(msg)
	window.advPool = nil
	window.await = make(chan struct{})
	window.cancel = nil
	window.set = true
	window.startIndex = 0
	window.endIndex = 0
	window.err = nil
}

func setError(err error) {
	window.mu.Lock()
	window.err = err
	window.mu.Unlock()
}

func newAdvertisingPacket(msgData []byte, sequenceNum uint8) *bluetooth.Advertisement {
	data := make([]byte, 1, spec.ByteSizeManufacturerData)
	data[0] = sequenceNum
	data = append(data, msgData...)

	adapter.mu.Lock()
	adv := adapter.adapter.NewAdvertisement()
	adapter.mu.Unlock()
	advOpns := bluetooth.AdvertisementOptions{
		ManufacturerData: []bluetooth.ManufacturerDataElement{
			{
				CompanyID: 0xff,
				Data:      data,
			},
		},
	}

	adv.Configure(advOpns)
	return adv
}

func advertiseBlockLocked(start int, end int) error {
	var err error
	for index := start; index < end; index++ {
		slot := window.advPool[index%len(window.advPool)]
		if !slot.Started() {
			err = slot.Start()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func stopAdvertisingBlockLocked(start int, end int) error {
	var err error
	fmt.Printf("start index: %d, end index: %d\n", start, end)
	for index := start; index < end; index++ {
		slot := window.advPool[index%len(window.advPool)]
		if !slot.Started() {
			fmt.Println("packet not advertising: ", slot.Started())
			continue
		}
		err = slot.Stop()
		if err != nil {
			return err
		}
	}
	return nil
}

func endAdvertisement() {
	window.mu.Lock()
	defer window.mu.Unlock()
	window.set = false
	err := stopAdvertisingBlockLocked(window.startIndex, window.endIndex)
	if err != nil {
		fmt.Printf("error while ending advertisement, err: %+v\n", err)
	}
}

func slideRight() error {

	window.mu.Lock()
	defer window.mu.Unlock()
	slotsBefore, _ := adapter.adapter.GetAdvertisementSlots()

	err := stopAdvertisingBlockLocked(window.startIndex, window.startIndex+1)
	if err != nil {
		return err
	}

	slotsAfter, _ := adapter.adapter.GetAdvertisementSlots()
	if slotsAfter <= slotsBefore {
		fmt.Println("advertisements slot was not released")
	}

	window.startIndex++
	err = advertiseBlockLocked(window.endIndex, window.endIndex+1)
	if err != nil {
		return err
	}
	window.endIndex++
	return nil
}

func slidingWindow(ctx context.Context) {

	var skipSliding bool

	payloadSize := spec.ByteSizeManufacturerData - 1
	window.mu.Lock()
	numPackets := (len(window.data) + payloadSize - 1) / payloadSize
	window.advPool = make([]*bluetooth.Advertisement, 0, numPackets)

	window.endIndex = int(advertisementSlots)
	window.startIndex = 0
	if numPackets <= int(advertisementSlots) {
		skipSliding = true
		window.endIndex = numPackets
	}

	newPacketStart := 0
	var sequenceNum uint8
	for range numPackets {
		newPacketEnd := min(len(window.data), newPacketStart+payloadSize)
		window.advPool = append(window.advPool, newAdvertisingPacket(window.data[newPacketStart:newPacketEnd], sequenceNum))
		newPacketStart = newPacketEnd
		sequenceNum++
	}

	// start advertising the window
	err := advertiseBlockLocked(window.startIndex, window.endIndex)
	if err != nil {
		window.err = err
		defer window.mu.Unlock() // whatever you do, unlock this mutex before you go
		err1 := stopAdvertisingBlockLocked(window.startIndex, window.endIndex)
		if err1 != nil {
			fmt.Printf("failed to start advertising, error: %+v; failed to stop advertising, error: %+v\n", err, err1)
			return
		}
		fmt.Printf("failed to start advertising, error: %+v\n", err)
		return
	}

	window.mu.Unlock()

	if skipSliding {
		return
	}

	for {
		select {
		case <-time.After(time.Second * 2):
			err = slideRight()
			if err != nil {
				fmt.Printf("failed to slide window, error: %+v\n", err)
				setError(err)
				endAdvertisement()
				return
			}
		case <-ctx.Done():
			endAdvertisement()
			window.mu.Lock()
			window.await <- struct{}{}
			window.mu.Unlock()
			return
		}
	}

}

func FetchError() error {
	window.mu.Lock()
	defer window.mu.Unlock()
	return window.err
}

func Stop() {
	window.mu.Lock()
	if !window.set {
		window.mu.Unlock()
		return
	}
	window.mu.Unlock()
	endAdvertisement()
}

func Broadcast(c *BeaconCmd) spec.ErrorResponse {

	err := validateMessage(c.Message)
	if err != nil {
		return spec.ErrorResponse{
			Error:  spec.ErrInvalidMessage,
			ErrMsg: err.Error(),
		}
	}

	window.mu.Lock()
	defer window.mu.Unlock()
	if window.set {
		return spec.ErrorResponse{
			Error: spec.ErrBroadcastAlreadyInProgress,
		}
	}
	setWindowLocked(c.Message)

	err = setAdapter()
	if err != nil {
		return spec.ErrorResponse{
			Error:  spec.ErrFailedToEnableAdapter,
			ErrMsg: err.Error(),
		}
	}

	advertisementSlots, err = adapter.adapter.GetAdvertisementSlots()
	if err != nil {
		return spec.ErrorResponse{
			Error:  spec.ErrUnknown,
			ErrMsg: err.Error(),
		}
	}

	myAddress, _ := adapter.adapter.Address()
	fmt.Println("broadcasting from address: ", myAddress.String())

	var ctx context.Context
	ctx, window.cancel = context.WithCancel(context.Background())

	go slidingWindow(ctx)

	return spec.ErrorResponse{}
}
