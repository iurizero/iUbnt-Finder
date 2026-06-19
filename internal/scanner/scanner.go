package scanner

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"syscall"
	"time"
)

type udpConn interface {
	ReadFromUDP([]byte) (int, *net.UDPAddr, error)
	WriteToUDP([]byte, *net.UDPAddr) (int, error)
	SetDeadline(time.Time) error
	Close() error
	SyscallConn() (syscall.RawConn, error)
}

type Scanner struct {
	Config    Config
	Decoder   *Decoder
	LookupMAC func(string) string
	Now       func() time.Time
	OpenConn  func() (udpConn, error)
}

func NewScanner(timeout time.Duration) *Scanner {
	return &Scanner{
		Config: Config{
			BroadcastAddr: "255.255.255.255:10001",
			Timeout:       timeout,
		},
		Decoder:   NewDecoder(),
		LookupMAC: lookupMACFromARP,
		Now:       time.Now,
	}
}

func (s *Scanner) Scan() ([]Device, error) {
	targets, err := resolveTargets(s.Config)
	if err != nil {
		return nil, err
	}

	openConn := s.OpenConn
	if openConn == nil {
		openConn = func() (udpConn, error) {
			return net.ListenUDP("udp4", localUDPAddr(s.Config.LocalAddr))
		}
	}

	conn, err := openConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if err := enableBroadcast(conn); err != nil {
		return nil, err
	}

	now := s.Now
	if now == nil {
		now = time.Now
	}

	if err := conn.SetDeadline(now().Add(s.Config.Timeout)); err != nil {
		return nil, err
	}

	probe := s.Decoder.ProbePacket()
	for _, target := range targets {
		if _, err := conn.WriteToUDP(probe, target); err != nil {
			return nil, err
		}
	}

	seen := make(map[string]struct{})
	devices := make([]Device, 0)
	buffer := make([]byte, 4096)
	lookupMAC := s.LookupMAC
	if lookupMAC == nil {
		lookupMAC = lookupMACFromARP
	}

	for {
		n, remote, err := conn.ReadFromUDP(buffer)
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) || isTimeout(err) {
				break
			}
			return devices, err
		}

		decoded := s.Decoder.Parse(buffer[:n])
		if decoded == nil {
			continue
		}

		device := Device{
			IPAddress:    remote.IP.String(),
			MACAddress:   lookupMAC(remote.IP.String()),
			DeviceName:   decoded["device_name"],
			ProductModel: decoded["product_model"],
			ModelID:      decoded["model_id"],
			BuildVersion: decoded["build_version"],
			NetworkID:    decoded["network_id"],
		}

		key := device.IPAddress + "|" + device.ProductModel + "|" + device.DeviceName
		if _, ok := seen[key]; ok {
			continue
		}

		seen[key] = struct{}{}
		devices = append(devices, device)
	}

	return devices, nil
}

func resolveTargets(config Config) ([]*net.UDPAddr, error) {
	targetSpecs := config.Targets
	if len(targetSpecs) == 0 {
		targetSpecs = []string{config.BroadcastAddr}
	}

	targets := make([]*net.UDPAddr, 0, len(targetSpecs))
	for _, spec := range targetSpecs {
		addr, err := resolveTarget(spec)
		if err != nil {
			return nil, err
		}
		targets = append(targets, addr)
	}

	return targets, nil
}

func localUDPAddr(spec string) *net.UDPAddr {
	if strings.TrimSpace(spec) == "" {
		return &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	}

	ip := net.ParseIP(strings.TrimSpace(spec))
	if ip == nil {
		return &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	}

	return &net.UDPAddr{IP: ip, Port: 0}
}

func resolveTarget(spec string) (*net.UDPAddr, error) {
	trimmed := strings.TrimSpace(spec)
	if trimmed == "" {
		return nil, fmt.Errorf("target address is empty")
	}

	if _, _, err := net.SplitHostPort(trimmed); err != nil {
		if strings.Contains(err.Error(), "missing port in address") {
			trimmed = net.JoinHostPort(trimmed, "10001")
		} else {
			return nil, err
		}
	}

	return net.ResolveUDPAddr("udp4", trimmed)
}

func isTimeout(err error) bool {
	netErr, ok := err.(net.Error)
	return ok && netErr.Timeout()
}

func enableBroadcast(conn udpConn) error {
	raw, err := conn.SyscallConn()
	if err != nil {
		return err
	}

	var controlErr error
	if err := raw.Control(func(fd uintptr) {
		controlErr = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
	}); err != nil {
		return err
	}

	return controlErr
}

func FormatDeviceInfo(device Device) string {
	productModel := device.ProductModel
	if productModel == "" {
		productModel = "Desconhecido"
	}

	deviceName := device.DeviceName
	if deviceName == "" {
		deviceName = "N/A"
	}

	modelID := device.ModelID
	if modelID == "" {
		modelID = "N/A"
	}

	buildVersion := device.BuildVersion
	if buildVersion == "" {
		buildVersion = "N/A"
	}

	networkID := device.NetworkID
	if networkID == "" {
		networkID = "N/A"
	}

	macAddress := device.MACAddress
	if macAddress == "" {
		macAddress = "N/A"
	}

	return fmt.Sprintf(
		"\n=== %s ===\nIP: %s\nNome: %s\nModelo: %s\nVersão: %s\nRede: %s\nMAC: %s\n",
		productModel,
		emptyAsNA(device.IPAddress),
		deviceName,
		modelID,
		buildVersion,
		networkID,
		macAddress,
	)
}

func emptyAsNA(value string) string {
	if value == "" {
		return "N/A"
	}
	return value
}
