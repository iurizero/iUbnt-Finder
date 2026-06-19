package scanner

import (
	"bytes"
	"testing"
)

func TestDecoderProbePacketReturnsCopy(t *testing.T) {
	decoder := NewDecoder()

	first := decoder.ProbePacket()
	second := decoder.ProbePacket()

	if !bytes.Equal(first, second) {
		t.Fatalf("expected probe packets to match")
	}

	first[0] = 0x99
	third := decoder.ProbePacket()
	if third[0] != 0x01 {
		t.Fatalf("expected probe packet to be immutable copy, got %x", third[0])
	}
}

func TestDecoderParse(t *testing.T) {
	decoder := NewDecoder()
	packet := buildTestPacket(
		[]tlvField{
			{marker: 0x0b, value: "AP-Sala"},
			{marker: 0x14, value: "UAP-AC-PRO"},
			{marker: 0x0c, value: "UAP"},
			{marker: 0x03, value: "6.5.28"},
			{marker: 0x0d, value: "Rede-Corporativa"},
		},
	)

	got := decoder.Parse(packet)
	if got == nil {
		t.Fatalf("expected parsed data, got nil")
	}

	want := map[string]string{
		"device_name":   "AP-Sala",
		"product_model": "UAP-AC-PRO",
		"model_id":      "UAP",
		"build_version": "6.5.28",
		"network_id":    "Rede-Corporativa",
	}

	for key, wantValue := range want {
		if got[key] != wantValue {
			t.Fatalf("field %s: want %q, got %q", key, wantValue, got[key])
		}
	}
}

func TestDecoderParseRejectsInvalidPacket(t *testing.T) {
	decoder := NewDecoder()

	if got := decoder.Parse([]byte{0x02, 0x00, 0x00, 0x01}); got != nil {
		t.Fatalf("expected nil for invalid header, got %#v", got)
	}

	if got := decoder.Parse([]byte{0x01, 0x00, 0x00}); got != nil {
		t.Fatalf("expected nil for too-short packet, got %#v", got)
	}
}

type tlvField struct {
	marker byte
	value  string
}

func buildTestPacket(fields []tlvField) []byte {
	payload := make([]byte, 0)
	for _, field := range fields {
		payload = append(payload, field.marker)
		payload = append(payload, byte(len(field.value)>>8), byte(len(field.value)))
		payload = append(payload, []byte(field.value)...)
	}

	packet := []byte{0x01, 0x00, 0x00, byte(len(payload) + 1)}
	packet = append(packet, payload...)
	return packet
}
