package main

import (
	"flag"
	"log"
	"time"

	"github.com/yznts/elklib"
)

func main() {
	// Parse command line flags
	address := flag.String("address", "", "Bluetooth MAC address of the ELK device (required)")
	flag.Parse()

	if *address == "" {
		log.Fatal("Please provide device address with -address flag")
	}

	// Create device
	device := elklib.NewDevice(*address)

	// Connect
	log.Printf("Connecting to device at %s...", *address)
	if err := device.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer device.Disconnect()

	log.Println("Connected successfully!")

	// Power on
	log.Println("Powering on...")
	if err := device.PowerOn(); err != nil {
		log.Fatalf("Failed to power on: %v", err)
	}
	time.Sleep(500 * time.Millisecond)

	// Set red color
	log.Println("Setting color to red...")
	if err := device.SetColor(255, 0, 0); err != nil {
		log.Fatalf("Failed to set color: %v", err)
	}
	time.Sleep(2 * time.Second)

	// Set green color
	log.Println("Setting color to green...")
	if err := device.SetColor(0, 255, 0); err != nil {
		log.Fatalf("Failed to set color: %v", err)
	}
	time.Sleep(2 * time.Second)

	// Set blue color
	log.Println("Setting color to blue...")
	if err := device.SetColor(0, 0, 255); err != nil {
		log.Fatalf("Failed to set color: %v", err)
	}
	time.Sleep(2 * time.Second)

	// Set brightness to 50%
	log.Println("Setting brightness to 50%...")
	if err := device.SetBrightness(50); err != nil {
		log.Fatalf("Failed to set brightness: %v", err)
	}
	time.Sleep(2 * time.Second)

	// Set rainbow crossfade effect
	log.Println("Setting rainbow crossfade effect...")
	if err := device.SetEffect(elklib.EffectCrossfadeRGBYCMW); err != nil {
		log.Fatalf("Failed to set effect: %v", err)
	}
	time.Sleep(2 * time.Second)

	// Set effect speed to 80%
	log.Println("Setting effect speed to 80%...")
	if err := device.SetEffectSpeed(80); err != nil {
		log.Fatalf("Failed to set effect speed: %v", err)
	}

	log.Println("Demo complete! Effect will continue running.")
	time.Sleep(5 * time.Second)

	// Power off
	log.Println("Powering off...")
	if err := device.PowerOff(); err != nil {
		log.Fatalf("Failed to power off: %v", err)
	}

	log.Println("Done!")
}
