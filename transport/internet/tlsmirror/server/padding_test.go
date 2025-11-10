package server

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPack(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	paddingLength := 3

	packed := Pack(data, paddingLength)

	// Expected format: data (5 bytes) + padding (3 bytes) + data length (4 bytes) = 12 bytes
	expectedLength := len(data) + paddingLength + 4
	if len(packed) != expectedLength {
		t.Errorf("Expected packed length %d, got %d", expectedLength, len(packed))
	}

	// Verify data is at the beginning
	if diff := cmp.Diff(packed[:5], data); diff != "" {
		t.Errorf("Data mismatch (-got +want):\n%s", diff)
	}

	// Verify padding (zeros)
	expectedPadding := []byte{0x00, 0x00, 0x00}
	if diff := cmp.Diff(packed[5:8], expectedPadding); diff != "" {
		t.Errorf("Padding mismatch (-got +want):\n%s", diff)
	}

	// Verify length encoding (big-endian uint32 = 5)
	expectedLengthEncoding := []byte{0x00, 0x00, 0x00, 0x05}
	if diff := cmp.Diff(packed[8:], expectedLengthEncoding); diff != "" {
		t.Errorf("Length encoding mismatch (-got +want):\n%s", diff)
	}
}

func TestPackWithZeroPadding(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03}
	paddingLength := 0

	packed := Pack(data, paddingLength)

	// Expected format: data (3 bytes) + data length (4 bytes) = 7 bytes
	expectedLength := len(data) + 4
	if len(packed) != expectedLength {
		t.Errorf("Expected packed length %d, got %d", expectedLength, len(packed))
	}

	// Verify data
	if diff := cmp.Diff(packed[:3], data); diff != "" {
		t.Errorf("Data mismatch (-got +want):\n%s", diff)
	}

	// Verify length encoding
	expectedLengthEncoding := []byte{0x00, 0x00, 0x00, 0x03}
	if diff := cmp.Diff(packed[3:], expectedLengthEncoding); diff != "" {
		t.Errorf("Length encoding mismatch (-got +want):\n%s", diff)
	}
}

func TestPackWithLargePadding(t *testing.T) {
	data := []byte{0xAA, 0xBB}
	paddingLength := 100

	packed := Pack(data, paddingLength)

	expectedLength := len(data) + paddingLength + 4
	if len(packed) != expectedLength {
		t.Errorf("Expected packed length %d, got %d", expectedLength, len(packed))
	}

	// Verify data
	if diff := cmp.Diff(packed[:2], data); diff != "" {
		t.Errorf("Data mismatch (-got +want):\n%s", diff)
	}

	// Verify all padding bytes are zero
	for i := 2; i < 102; i++ {
		if packed[i] != 0x00 {
			t.Errorf("Expected padding byte at position %d to be 0x00, got 0x%02x", i, packed[i])
		}
	}
}

func TestPad(t *testing.T) {
	tests := []struct {
		name           string
		paddingLength  int
		expectedLength int
	}{
		{"Zero padding", 0, 0},
		{"One byte", 1, 1},
		{"Two bytes", 2, 2},
		{"Three bytes", 3, 3},
		{"Four bytes", 4, 4},
		{"Five bytes", 5, 5},
		{"Large padding", 100, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			padding := Pad(tt.paddingLength)
			if len(padding) != tt.expectedLength {
				t.Errorf("Expected padding length %d, got %d", tt.expectedLength, len(padding))
			}

			// Verify all bytes are zero
			for i, b := range padding {
				if b != 0x00 {
					t.Errorf("Expected padding byte at position %d to be 0x00, got 0x%02x", i, b)
				}
			}
		})
	}
}

func TestPadNegative(t *testing.T) {
	padding := Pad(-1)
	if padding != nil {
		t.Error("Expected nil for negative padding length")
	}
}

func TestUnpack(t *testing.T) {
	// Create a packed message
	originalData := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	paddingLength := 3
	packed := Pack(originalData, paddingLength)

	// Unpack it
	data, extractedPaddingLength := Unpack(packed)

	// Verify data
	if diff := cmp.Diff(data, originalData); diff != "" {
		t.Errorf("Unpacked data mismatch (-got +want):\n%s", diff)
	}

	// Verify padding length
	if extractedPaddingLength != paddingLength {
		t.Errorf("Expected padding length %d, got %d", paddingLength, extractedPaddingLength)
	}
}

func TestUnpackZeroPadding(t *testing.T) {
	originalData := []byte{0xAA, 0xBB, 0xCC}
	paddingLength := 0
	packed := Pack(originalData, paddingLength)

	data, extractedPaddingLength := Unpack(packed)

	if diff := cmp.Diff(data, originalData); diff != "" {
		t.Errorf("Unpacked data mismatch (-got +want):\n%s", diff)
	}

	if extractedPaddingLength != paddingLength {
		t.Errorf("Expected padding length %d, got %d", paddingLength, extractedPaddingLength)
	}
}

func TestUnpackTooShort(t *testing.T) {
	// Data shorter than 4 bytes
	shortData := []byte{0x01, 0x02, 0x03}

	data, paddingLength := Unpack(shortData)

	if data != nil {
		t.Error("Expected nil data for short packet")
	}

	if paddingLength != len(shortData) {
		t.Errorf("Expected padding length %d, got %d", len(shortData), paddingLength)
	}
}

func TestUnpackInvalidLength(t *testing.T) {
	// Create invalid data with length field > actual data size
	invalidData := []byte{
		0x01, 0x02, 0x03, // 3 bytes of data
		0x00, 0x00, 0x00, 0xFF, // length field claiming 255 bytes
	}

	data, paddingLength := Unpack(invalidData)

	if data != nil {
		t.Error("Expected nil data for invalid length field")
	}

	if paddingLength != 0 {
		t.Errorf("Expected padding length 0, got %d", paddingLength)
	}
}

func TestPackUnpackRoundTrip(t *testing.T) {
	testCases := []struct {
		name          string
		data          []byte
		paddingLength int
	}{
		{"Small data, small padding", []byte{0x01, 0x02}, 2},
		{"Medium data, medium padding", []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE}, 10},
		{"Large data, large padding", make([]byte, 100), 50},
		{"Empty data, small padding", []byte{}, 5},
		{"Data with zero padding", []byte{0xFF, 0xFE, 0xFD}, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize test data
			for i := range tc.data {
				tc.data[i] = byte(i % 256)
			}

			// Pack and unpack
			packed := Pack(tc.data, tc.paddingLength)
			unpacked, extractedPaddingLength := Unpack(packed)

			// Verify data
			if diff := cmp.Diff(unpacked, tc.data); diff != "" {
				t.Errorf("Round trip data mismatch (-got +want):\n%s", diff)
			}

			// Verify padding length
			if extractedPaddingLength != tc.paddingLength {
				t.Errorf("Expected padding length %d, got %d", tc.paddingLength, extractedPaddingLength)
			}
		})
	}
}

func TestUnpackPurelyPaddedPacket(t *testing.T) {
	// Test packets that are only padding (< 4 bytes)
	testCases := []struct {
		name          string
		paddingLength int
	}{
		{"1 byte padding", 1},
		{"2 bytes padding", 2},
		{"3 bytes padding", 3},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			padding := Pad(tc.paddingLength)
			data, extractedPaddingLength := Unpack(padding)

			if data != nil {
				t.Error("Expected nil data for pure padding packet")
			}

			if extractedPaddingLength != tc.paddingLength {
				t.Errorf("Expected padding length %d, got %d", tc.paddingLength, extractedPaddingLength)
			}
		})
	}
}
