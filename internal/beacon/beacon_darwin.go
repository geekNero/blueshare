//go:build darwin

package beacon

import "fmt"

func Broadcast(_ *Message) error {
	return fmt.Errorf("broadcasting is not implemented for MacOS")
}
