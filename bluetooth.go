package elklib

import (
	"errors"
	"fmt"
	"sync"

	"tinygo.org/x/bluetooth"
)

const (
	// Service and characteristic UUIDs for ELK-BLE devices
	WriteUUID = "0000fff3-0000-1000-8000-00805f9b34fb"
	ReadUUID  = "0000fff4-0000-1000-8000-00805f9b34fb"
)

var (
	Adapter        = bluetooth.DefaultAdapter
	AdapterEnabled = false

	// connectMu serializes Enable + Connect calls — CoreBluetooth cannot handle
	// concurrent connection attempts from multiple goroutines.
	connectMu sync.Mutex

	ErrNotConnected   = errors.New("device not connected")
	ErrCharacteristic = errors.New("characteristic not found")
)

// Device represents a Bluetooth ELK LED device
type Device struct {
	adapter   *bluetooth.Adapter
	device    *bluetooth.Device
	writeChar bluetooth.DeviceCharacteristic
	connected bool
	address   string
}

// NewDevice creates a new device instance with manual address (MAC or UUID)
func NewDevice(address string) *Device {
	return &Device{
		address: address,
	}
}

// Connect connects to the device using the provided address (MAC address or UUID)
func (d *Device) Connect(params ...bluetooth.ConnectionParams) error {
	// Default params
	if len(params) == 0 {
		params = append(params, bluetooth.ConnectionParams{})
	}

	// Serialize Enable + Connect: concurrent attempts cause CoreBluetooth timeouts.
	connectMu.Lock()
	if !AdapterEnabled {
		if err := Adapter.Enable(); err != nil {
			connectMu.Unlock()
			return fmt.Errorf("failed to enable adapter: %w", err)
		}
		AdapterEnabled = true
	}
	d.adapter = Adapter

	addr, err := ParseAddress(d.address)
	if err != nil {
		connectMu.Unlock()
		return fmt.Errorf("invalid address: %w", err)
	}

	device, err := d.adapter.Connect(addr, params[0])
	connectMu.Unlock() // release before service discovery — that part is per-device safe
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	d.device = &device

	// Discover services
	services, err := device.DiscoverServices(nil)
	if err != nil {
		return fmt.Errorf("failed to discover services: %w", err)
	}

	// Find write characteristic
	found := false
	for _, svc := range services {
		chars, err := svc.DiscoverCharacteristics(nil)
		if err != nil {
			continue
		}

		for _, char := range chars {
			if char.UUID().String() == WriteUUID {
				d.writeChar = char
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		return ErrCharacteristic
	}

	d.connected = true
	return nil
}

// Disconnect disconnects from the device
func (d *Device) Disconnect() error {
	if d.device != nil && d.connected {
		err := d.device.Disconnect()
		d.connected = false
		return err
	}
	return nil
}

// sendCommand sends a command to the device
func (d *Device) sendCommand(cmd []byte) error {
	if !d.connected {
		return ErrNotConnected
	}

	_, err := d.writeChar.WriteWithoutResponse(cmd)
	if err != nil {
		return fmt.Errorf("failed to write command: %w", err)
	}

	return nil
}

// PowerOn turns the device on
func (d *Device) PowerOn() error {
	cmd := []byte{0x7e, 0x00, 0x04, 0xf0, 0x00, 0x01, 0xff, 0x00, 0xef}
	return d.sendCommand(cmd)
}

// PowerOff turns the device off
func (d *Device) PowerOff() error {
	cmd := []byte{0x7e, 0x00, 0x04, 0x00, 0x00, 0x00, 0xff, 0x00, 0xef}
	return d.sendCommand(cmd)
}

// SetColor sets the RGB color
func (d *Device) SetColor(r, g, b uint8) error {
	cmd := []byte{0x7e, 0x00, 0x05, 0x03, r, g, b, 0x00, 0xef}
	return d.sendCommand(cmd)
}

// SetBrightness sets the brightness (0-100)
func (d *Device) SetBrightness(value uint8) error {
	if value > 100 {
		value = 100
	}
	cmd := []byte{0x7e, 0x00, 0x01, value, 0x00, 0x00, 0x00, 0x00, 0xef}
	return d.sendCommand(cmd)
}

// SetEffect sets a light effect mode
func (d *Device) SetEffect(effectCode uint8) error {
	cmd := []byte{0x7e, 0x00, 0x03, effectCode, 0x03, 0x00, 0x00, 0x00, 0xef}
	return d.sendCommand(cmd)
}

// SetEffectSpeed sets the speed of the current effect (0-100)
func (d *Device) SetEffectSpeed(value uint8) error {
	if value > 100 {
		value = 100
	}
	cmd := []byte{0x7e, 0x00, 0x02, value, 0x00, 0x00, 0x00, 0x00, 0xef}
	return d.sendCommand(cmd)
}

// Ping checks whether the BLE connection is still alive by querying the
// CoreBluetooth peripheral state. No packets are sent to the device.
// Returns ErrNotConnected if the connection has been lost.
func (d *Device) Ping() error {
	if d.device == nil {
		return ErrNotConnected
	}
	ok, err := d.device.Connected()
	if err != nil {
		return fmt.Errorf("connection check: %w", err)
	}
	if !ok {
		return ErrNotConnected
	}
	return nil
}
