package main

import (
	"time"

	"github.com/soypat/cyw43439"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

var (
	serviceUUID = [16]byte{0xa0, 0xb4, 0x00, 0x01, 0x92, 0x6d, 0x4d, 0x61, 0x98, 0xdf, 0x8c, 0x5c, 0x62, 0xee, 0x53, 0xb3}
	charUUID    = [16]byte{0xa0, 0xb4, 0x00, 0x02, 0x92, 0x6d, 0x4d, 0x61, 0x98, 0xdf, 0x8c, 0x5c, 0x62, 0xee, 0x53, 0xb3}
)

func main() {
	time.Sleep(time.Second)

	// The following code is required for ensuring the LED is available.
	// - Placing it after bluetooth.DefaultAdapter.Enable() prevents the
	//   bluetooth service being enabled.
	// - Placing it before bluetooth.DefaultAdapter.Enable() prevents the
	//   LED from functioning with a failure in the chan range at the
	//   bottom of "gpio set: pollForIoctl timeout".
	//
	dev := cyw43439.NewPicoWDevice()
	must("initialise device", dev.Init(cyw43439.DefaultWifiBluetoothConfig()))

	println("starting")
	adapter.Use(dev)

	adv := adapter.DefaultAdvertisement()
	must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: "LED flasher",
	}))
	must("start adv", adv.Start())
	led := make(chan bool)
	var ledData [1]byte
	var ledCharacteristic bluetooth.Characteristic
	must("add service", adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.NewUUID(serviceUUID),
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: &ledCharacteristic,
				UUID:   bluetooth.NewUUID(charUUID),
				Value:  ledData[:],
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission | bluetooth.CharacteristicWriteWithoutResponsePermission,
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					if offset != 0 || len(value) != 1 {
						return
					}
					ledData[0] = value[0]
					select {
					case led <- value[0] != 0:
					default:
					}
				},
			},
		},
	}))

	for on := range led {
		println(on)
		err := dev.GPIOSet(0, on)
		if err != nil {
			println("gpio set: " + err.Error())
		}
	}
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
