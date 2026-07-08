package scan

import (
	"blueshare/internal/spec"
	"blueshare/internal/utility"
	"fmt"
	"net/rpc"
)

type ScanCmd struct {
	Fetch    bool `help:"Fetch currently scanned message"`
	StopScan bool `help:"Stop scanning"`
}

func (c *ScanCmd) Run() error {
	socketAddr := utility.GetSocketPath(spec.SocketAddr)
	client, err := rpc.Dial("unix", socketAddr)
	if err != nil {
		return fmt.Errorf("error while connecting to socket, ensure blueshare daemon is running, error: %+v", err)
	}

	defer client.Close()

	if c.Fetch {
		var reply string
		err = client.Call("Ritual.FetchMessage", struct{}{}, &reply)
		if err != nil {
			return fmt.Errorf("error while calling daemon method, error: %+v", err)
		}

		fmt.Println("message found is: ", reply)
		return nil
	} else if c.StopScan {
		var reply error
		err = client.Call("Ritual.StopScan", struct{}{}, &reply)
		if err != nil {
			return fmt.Errorf("error while calling daemon method, error: %+v", err)
		}
		if reply != nil {
			return fmt.Errorf("failed to stop scanning, err: %+v", reply)
		}
	} else {
		var reply spec.ErrorResponse
		err = client.Call("Ritual.Scan", struct{}{}, &reply)
		if err != nil {
			return fmt.Errorf("error while calling daemon method, error: %+v", err)
		}
		if reply.Err != nil {
			return fmt.Errorf("%s, err: %+v", spec.ErrorMap[reply.Error], reply.Err)
		}
		if reply.Error != spec.Error(0) {
			fmt.Println(spec.ErrorMap[reply.Error])
		}
	}
	return nil
}
