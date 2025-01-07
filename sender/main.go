package main

import (
	"fmt"
	"strconv"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

var (
	serviceUUID = [16]byte{0xa0, 0xb4, 0x00, 0x01, 0x92, 0x6d, 0x4d, 0x61, 0x98, 0xdf, 0x8c, 0x5c, 0x62, 0xee, 0x53, 0xb3}
	charUUID    = [16]byte{0xa0, 0xb4, 0x00, 0x02, 0x92, 0x6d, 0x4d, 0x61, 0x98, 0xdf, 0x8c, 0x5c, 0x62, 0xee, 0x53, 0xb3}
)

func main() {
	err := adapter.Enable()
	if err != nil {
		fmt.Printf("failed to enable bluetooth: %v", err)
	}

	fmt.Println("scanning...")
	var dev bluetooth.Device
	err = adapter.Scan(func(adapter *bluetooth.Adapter, found bluetooth.ScanResult) {
		fmt.Printf("%s\r", found.Address)
		if found.LocalName() == "LED flasher" {
			fmt.Println("found device:", found.Address, found.RSSI, found.LocalName(), found.AdvertisementPayload.LocalName(), found.AdvertisementPayload.Bytes())
			dev, err = adapter.Connect(found.Address, bluetooth.ConnectionParams{})
			if err != nil {
				fmt.Printf("failed to connect: %v", err)
				return
			}
			adapter.StopScan()
		}
	})

	srvs, err := dev.DiscoverServices([]bluetooth.UUID{bluetooth.NewUUID(serviceUUID)})
	if err != nil {
		fmt.Printf("failed to discover service: %v", err)
		return
	}
	for _, srv := range srvs {
		fmt.Println(srv.UUID(), srv.String())
		chars, err := srv.DiscoverCharacteristics([]bluetooth.UUID{bluetooth.NewUUID(charUUID)})
		if err != nil {
			fmt.Printf("failed to discover characteristic: %v", err)
		}
		buf := make([]byte, 255)
		for _, char := range chars {
			n, err := char.Read(buf)
			if err != nil {
				fmt.Printf("failed to read characteristic: %v", err)
				continue
			}
			buf[0] = 1 - buf[0]
			fmt.Println("    data bytes", strconv.Itoa(n))
			fmt.Printf("    value = %v\n", buf[:n])
			n, err = char.WriteWithoutResponse(buf[:n])
			if err != nil {
				fmt.Printf("failed to write characteristic: %v", err)
			}
		}
	}
	err = dev.Disconnect()
	if err != nil {
		println(err)
	}
}
