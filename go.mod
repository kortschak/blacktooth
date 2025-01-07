module github.com/kortschak/blacktooth

go 1.23.4

// Cannot require tinygo.org/x/bluetooth v0.10.1-dev as this
// breaks bluetooth service discovery. Bisected to:
//
// commit c74fb68742be8788734b1f65830aa75289e2408e
// Author: deadprogram <ron@hybridgroup.com>
// Date:   Sun Nov 17 09:35:50 2024 +0100
// 
//     hci: contains several corrections needed for HCI on both ninafw and cyw43439
//     
//     * separate read and write buffers to accomodate read buffer with pending partial data to be handled
//     * only read when there is data pending to avoid putting cwy43439 into forever loop
//     * handle any data in buffer from read that contains data from more than 1 packet
//     * skip past data for sync packets since we do not handle them anyhow
//     * some better naming and comments to hopefully clarify a bit what is going on
//     
//     Signed-off-by: deadprogram <ron@hybridgroup.com>
// 
// :100644 100644 5a4edffc7c850d8f580cf0b6e706d55d1e4e7d89 c0b1148b78379ac42eadcd4f907b9c69f0ec6a44 M	hci.go
//
// Confirmed by reverting this commit.
require tinygo.org/x/bluetooth v0.10.0

require (
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/saltosystems/winrt-go v0.0.0-20240509164145-4f7860a3bd2b // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/soypat/seqs v0.0.0-20240527012110-1201bab640ef // indirect
	github.com/tinygo-org/cbgo v0.0.4 // indirect
	github.com/tinygo-org/pio v0.0.0-20231216154340-cd888eb58899 // indirect
	golang.org/x/exp v0.0.0-20230728194245-b0cb94b80691 // indirect
	golang.org/x/sys v0.11.0 // indirect

	// Necessary to pick up ioctl fixes in cyw43439.
	github.com/soypat/cyw43439 v0.0.0-20250106095300-90bf0c1db251 // indirect
)
