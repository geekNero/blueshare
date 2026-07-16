package main

import (
	"blueshare/internal/beacon"
	"blueshare/internal/scan"
	"blueshare/internal/spec"
	"blueshare/internal/utility"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
)

type Ritual struct {
}

func (r *Ritual) Broadcast(c beacon.BeaconCmd, err *spec.ErrorResponse) error {
	*err = beacon.Broadcast(&c)
	return nil
}

func (r *Ritual) StopBroadcast(_ Ritual, _ *Ritual) error {
	beacon.Stop()
	return nil
}

func (r *Ritual) Scan(_ Ritual, err *spec.ErrorResponse) error {
	*err = scan.StartScan()
	return nil
}

func (r *Ritual) StopScan(_ Ritual, err *string) error {
	e := scan.StopScan()
	if e != nil {
		*err = e.Error()
	} else {
		*err = ""
	}
	return nil
}

func (r *Ritual) FetchMessage(_ Ritual, messages *[]string) error {
	messages = scan.FetchMessages()
	return nil
}

func main() {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	socketAddr := utility.GetSocketPath(spec.SocketAddr)
	os.Remove(socketAddr)

	ritual := new(Ritual)
	rpc.Register(ritual)

	listener, err := net.Listen("unix", socketAddr)
	if err != nil {
		panic(err)
	}

	go func() {

		fmt.Println("Daemon listening on socket ...")
		rpc.Accept(listener)
	}()

	<-sigs // Blocks until systemd or user sends a kill signal

	fmt.Println("\nSignal received. Purging socket and returning to the void...")
	if fetchErr := beacon.FetchError(); fetchErr != nil {
		fmt.Printf("Error found in broadcasting, err: %+v", fetchErr)
	}
	if err := scan.StopScan(); err != nil {
		fmt.Println(err)
	}
	listener.Close()
	os.Remove(socketAddr)
	os.Exit(0)

}
