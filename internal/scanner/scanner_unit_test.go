package scanner

import (
	"errors"
	"net"
	"sync"
	"syscall"
	"testing"
	"time"
)

func TestScannerScanWithFakeConnection(t *testing.T) {
	fake := &fakeUDPConn{
		readPackets: [][]byte{
			buildTestPacket([]tlvField{
				{marker: 0x0b, value: "AP-Sala"},
				{marker: 0x14, value: "UAP-AC-PRO"},
				{marker: 0x0c, value: "UAP"},
				{marker: 0x03, value: "6.5.28"},
				{marker: 0x0d, value: "Rede-Corporativa"},
			}),
			buildTestPacket([]tlvField{
				{marker: 0x0b, value: "AP-Sala"},
				{marker: 0x14, value: "UAP-AC-PRO"},
				{marker: 0x0c, value: "UAP"},
				{marker: 0x03, value: "6.5.28"},
				{marker: 0x0d, value: "Rede-Corporativa"},
			}),
		},
		remote: &net.UDPAddr{IP: net.ParseIP("192.168.1.100"), Port: 10001},
		err:    errTimeout{},
	}

	sc := NewScanner(5 * time.Second)
	sc.Config.BroadcastAddr = "127.0.0.1:10001"
	sc.LookupMAC = func(string) string { return "AA:BB:CC:DD:EE:FF" }
	sc.Now = func() time.Time { return time.Unix(1700000000, 0) }
	sc.OpenConn = func() (udpConn, error) { return fake, nil }

	devices, err := sc.Scan()
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	if len(devices) != 1 {
		t.Fatalf("expected 1 device, got %d", len(devices))
	}

	device := devices[0]
	if device.IPAddress != "192.168.1.100" {
		t.Fatalf("expected IP 192.168.1.100, got %s", device.IPAddress)
	}
	if device.MACAddress != "AA:BB:CC:DD:EE:FF" {
		t.Fatalf("expected injected MAC, got %s", device.MACAddress)
	}
	if device.DeviceName != "AP-Sala" || device.ProductModel != "UAP-AC-PRO" {
		t.Fatalf("unexpected device data: %#v", device)
	}

	if fake.deadline.IsZero() {
		t.Fatal("expected deadline to be set")
	}

	wantDeadline := time.Unix(1700000000, 0).Add(5 * time.Second)
	if !fake.deadline.Equal(wantDeadline) {
		t.Fatalf("unexpected deadline: want %v got %v", wantDeadline, fake.deadline)
	}

	if len(fake.writes) != 1 {
		t.Fatalf("expected 1 probe write, got %d", len(fake.writes))
	}
	if string(fake.writes[0]) != string(NewDecoder().ProbePacket()) {
		t.Fatalf("unexpected probe payload: %v", fake.writes[0])
	}

	if fake.closeCount != 1 {
		t.Fatalf("expected connection to be closed once, got %d", fake.closeCount)
	}
}

type errTimeout struct{}

func (errTimeout) Error() string   { return "timeout" }
func (errTimeout) Timeout() bool   { return true }
func (errTimeout) Temporary() bool { return true }

type fakeUDPConn struct {
	mu          sync.Mutex
	readPackets [][]byte
	remote      *net.UDPAddr
	err         error
	readIndex   int
	writes      [][]byte
	deadline    time.Time
	closeCount  int
}

func (f *fakeUDPConn) ReadFromUDP(p []byte) (int, *net.UDPAddr, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.readIndex >= len(f.readPackets) {
		return 0, nil, f.err
	}

	packet := f.readPackets[f.readIndex]
	f.readIndex++
	copy(p, packet)

	if f.readIndex >= len(f.readPackets) {
		return len(packet), f.remote, f.err
	}

	return len(packet), f.remote, nil
}

func (f *fakeUDPConn) WriteToUDP(p []byte, _ *net.UDPAddr) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	copied := append([]byte(nil), p...)
	f.writes = append(f.writes, copied)
	return len(p), nil
}

func (f *fakeUDPConn) SetDeadline(t time.Time) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.deadline = t
	return nil
}

func (f *fakeUDPConn) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.closeCount++
	return nil
}

func (f *fakeUDPConn) SyscallConn() (syscall.RawConn, error) {
	return fakeRawConn{}, nil
}

type fakeRawConn struct{}

func (fakeRawConn) Control(func(fd uintptr)) error    { return nil }
func (fakeRawConn) Read(func(fd uintptr) bool) error  { return nil }
func (fakeRawConn) Write(func(fd uintptr) bool) error { return nil }

var _ udpConn = (*fakeUDPConn)(nil)
var _ syscall.RawConn = fakeRawConn{}

func TestScannerScanWithOpenConnError(t *testing.T) {
	sc := NewScanner(time.Second)
	sc.OpenConn = func() (udpConn, error) {
		return nil, errors.New("boom")
	}

	if _, err := sc.Scan(); err == nil {
		t.Fatal("expected error")
	}
}

func TestResolveTargetAddsDefaultPort(t *testing.T) {
	addr, err := resolveTarget("192.168.1.10")
	if err != nil {
		t.Fatalf("resolve target: %v", err)
	}

	if got, want := addr.String(), "192.168.1.10:10001"; got != want {
		t.Fatalf("unexpected target: want %s got %s", want, got)
	}
}

func TestResolveTargetsUsesExplicitList(t *testing.T) {
	targets, err := resolveTargets(Config{
		BroadcastAddr: "255.255.255.255:10001",
		Targets:       []string{"192.168.1.10", "192.168.1.11:2000"},
	})
	if err != nil {
		t.Fatalf("resolve targets: %v", err)
	}

	if got, want := len(targets), 2; got != want {
		t.Fatalf("unexpected target count: want %d got %d", want, got)
	}

	if got, want := targets[0].String(), "192.168.1.10:10001"; got != want {
		t.Fatalf("unexpected first target: want %s got %s", want, got)
	}

	if got, want := targets[1].String(), "192.168.1.11:2000"; got != want {
		t.Fatalf("unexpected second target: want %s got %s", want, got)
	}
}

func TestLocalUDPAddr(t *testing.T) {
	if got := localUDPAddr(""); got.IP.String() != "0.0.0.0" || got.Port != 0 {
		t.Fatalf("unexpected default local addr: %#v", got)
	}

	if got := localUDPAddr("10.0.0.10"); got.IP.String() != "10.0.0.10" || got.Port != 0 {
		t.Fatalf("unexpected bound local addr: %#v", got)
	}
}
