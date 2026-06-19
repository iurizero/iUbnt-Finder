package scanner

import "time"

type Config struct {
	BroadcastAddr string
	Targets       []string
	LocalAddr     string
	Timeout       time.Duration
}

type Device struct {
	IPAddress    string
	MACAddress   string
	DeviceName   string
	ProductModel string
	ModelID      string
	BuildVersion string
	NetworkID    string
}
