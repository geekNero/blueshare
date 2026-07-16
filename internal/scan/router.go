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
		var replies []string
		err = client.Call("Ritual.FetchMessage", struct{}{}, &replies)
		if err != nil {
			return fmt.Errorf("error while calling daemon method, error: %+v", err)
		}

		fmt.Println("list of messages found:")
		for index, reply := range replies {
			fmt.Printf("%d: %s", index, reply)
		}
		return nil
	} else if c.StopScan {
		var reply string
		err = client.Call("Ritual.StopScan", struct{}{}, &reply)
		if err != nil {
			return fmt.Errorf("error while calling daemon method, error: %+v", err)
		}
		if reply != "" {
			return fmt.Errorf("failed to stop scanning, err: %s", reply)
		}
	} else {
		var reply spec.ErrorResponse
		err = client.Call("Ritual.Scan", struct{}{}, &reply)
		if err != nil {
			return fmt.Errorf("error while calling daemon method, error: %+v", err)
		}
		if reply.ErrMsg != "" {
			return fmt.Errorf("%s, err: %s", spec.ErrorMap[reply.Error], reply.ErrMsg)
		}
		if reply.Error != spec.Error(0) {
			fmt.Println(spec.ErrorMap[reply.Error])
		}
	}
	return nil
}
