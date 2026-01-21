package elklib

import "tinygo.org/x/bluetooth"

func ParseAddress(addr string) (bluetooth.Address, error) {
	uuid, err := bluetooth.ParseUUID(addr)
	if err == nil {
		return bluetooth.Address{UUID: uuid}, nil
	}
	return bluetooth.Address{}, err
}
