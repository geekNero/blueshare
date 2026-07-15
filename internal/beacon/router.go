package beacon

import (
	"blueshare/internal/spec"
	"blueshare/internal/utility"
	"fmt"
	"net/rpc"
)

type BeaconCmd struct {
	Message   string `default:"" name:"msg" short:"m" help:"String to be broadcast."`
	Frequency int    `default:"2" name:"frequency" short:"f" help:"Frequency at which the message should be broardcast.\nWithin range 1-3"`
	Once      bool   `help:"Only broadcast message once."`
	Stop      bool   `help:"Stop broadcasting"`
}

func (c *BeaconCmd) Run() error {
	// spec.CustomUUID, _ = bluetooth.ParseUUID(spec.CustomUUIDString)
	socketAddr := utility.GetSocketPath(spec.SocketAddr)
	client, err := rpc.Dial("unix", socketAddr)
	if err != nil {
		return fmt.Errorf("error while connecting to socket, ensure blueshare daemon is running, error: %+v", err)
	}

	defer client.Close()

	if c.Stop {
		err = client.Call("Ritual.StopBroadcast", struct{}{}, &struct{}{})
		if err != nil {
			return fmt.Errorf("error while calling daemon method, error: %+v", err)
		}
		return nil
	}

	var reply spec.ErrorResponse

	err = client.Call("Ritual.Broadcast", c, &reply)
	if err != nil {
		return fmt.Errorf("error while calling daemon method, error: %+v", err)
	}

	if reply.ErrMsg != "" {
		fmt.Printf("%s, err: %s\n", spec.ErrorMap[reply.Error], reply.ErrMsg)
	}
	return nil

}
