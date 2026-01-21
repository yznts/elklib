package elklib

import "tinygo.org/x/bluetooth"

func ParseAddress(addr string) (bluetooth.Address, error) {
	mac, err := bluetooth.ParseMAC(addr)
	if err == nil {
		return bluetooth.Address{MACAddress: bluetooth.MACAddress{MAC: mac}}, nil
	}
	return bluetooth.Address{}, err
}
