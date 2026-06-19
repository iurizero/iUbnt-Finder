package scanner

import (
	"bytes"
	"encoding/binary"
	"strings"
)

type Decoder struct {
	validHeader []byte
	probePacket []byte
	markers     map[byte]string
}

func NewDecoder() *Decoder {
	return &Decoder{
		validHeader: []byte{0x01, 0x00, 0x00},
		probePacket: []byte{0x01, 0x00, 0x00, 0x00},
		markers: map[byte]string{
			0x0b: "device_name",
			0x14: "product_model",
			0x0c: "model_id",
			0x03: "build_version",
			0x0d: "network_id",
		},
	}
}

func (d *Decoder) ProbePacket() []byte {
	return append([]byte(nil), d.probePacket...)
}

func (d *Decoder) Parse(data []byte) map[string]string {
	if len(data) < 4 || !bytes.Equal(data[:3], d.validHeader) {
		return nil
	}

	remaining := int(data[3]) - 1
	pos := 4
	result := make(map[string]string)

	for remaining > 0 {
		if pos >= len(data) {
			break
		}

		marker := data[pos]
		pos++
		remaining--

		if pos+2 > len(data) {
			break
		}

		size := int(binary.BigEndian.Uint16(data[pos : pos+2]))
		pos += 2
		remaining -= 2

		if pos+size > len(data) {
			break
		}

		value := strings.ToValidUTF8(string(data[pos:pos+size]), "")
		pos += size
		remaining -= size

		if key, ok := d.markers[marker]; ok {
			result[key] = value
		}
	}

	if len(result) == 0 {
		return nil
	}

	return result
}
