package scanner

import "testing"

func TestFormatDeviceInfo(t *testing.T) {
	device := Device{
		IPAddress:    "192.168.1.100",
		MACAddress:   "00:11:22:33:44:55",
		DeviceName:   "AP-Sala",
		ProductModel: "UAP-AC-PRO",
		ModelID:      "UAP",
		BuildVersion: "6.5.28",
		NetworkID:    "Rede-Corporativa",
	}

	got := FormatDeviceInfo(device)
	want := "\n=== UAP-AC-PRO ===\nIP: 192.168.1.100\nNome: AP-Sala\nModelo: UAP\nVersão: 6.5.28\nRede: Rede-Corporativa\nMAC: 00:11:22:33:44:55\n"

	if got != want {
		t.Fatalf("unexpected formatted output\nwant:\n%q\ngot:\n%q", want, got)
	}
}

func TestFormatDeviceInfoDefaults(t *testing.T) {
	got := FormatDeviceInfo(Device{})
	want := "\n=== Desconhecido ===\nIP: N/A\nNome: N/A\nModelo: N/A\nVersão: N/A\nRede: N/A\nMAC: N/A\n"

	if got != want {
		t.Fatalf("unexpected formatted output for defaults\nwant:\n%q\ngot:\n%q", want, got)
	}
}
