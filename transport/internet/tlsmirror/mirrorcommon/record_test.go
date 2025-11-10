package mirrorcommon

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/v2fly/v2ray-core/v5/transport/internet/tlsmirror"
)

type testPeeker struct {
	data []byte
}

func (p *testPeeker) Peek(n int) ([]byte, error) {
	if len(p.data) < n {
		return p.data, nil
	}
	return p.data[:n], nil
}

func TestPackTLSRecord(t *testing.T) {
	record := tlsmirror.TLSRecord{
		RecordType:            0x17,                // application_data
		LegacyProtocolVersion: [2]byte{0x03, 0x03}, // TLS 1.2
		RecordLength:          5,
		Fragment:              []byte{0x01, 0x02, 0x03, 0x04, 0x05},
	}

	packed := PackTLSRecord(record)

	expected := []byte{
		0x17,       // record type
		0x03, 0x03, // protocol version
		0x00, 0x05, // length
		0x01, 0x02, 0x03, 0x04, 0x05, // fragment
	}

	if diff := cmp.Diff(packed, expected); diff != "" {
		t.Errorf("PackTLSRecord mismatch (-got +want):\n%s", diff)
	}
}

func TestPeekTLSRecord(t *testing.T) {
	// Valid TLS record
	validRecord := []byte{
		0x16,       // handshake
		0x03, 0x03, // TLS 1.2
		0x00, 0x05, // length: 5
		0x01, 0x02, 0x03, 0x04, 0x05, // fragment
	}

	peeker := &testPeeker{data: validRecord}
	record, tryAgain, processed, err := PeekTLSRecord(peeker, nil)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if tryAgain != 0 {
		t.Errorf("Expected tryAgain to be 0, got %d", tryAgain)
	}

	if processed != 10 {
		t.Errorf("Expected processed to be 10, got %d", processed)
	}

	if record.RecordType != 0x16 {
		t.Errorf("Expected record type 0x16, got 0x%02x", record.RecordType)
	}

	if record.RecordLength != 5 {
		t.Errorf("Expected record length 5, got %d", record.RecordLength)
	}

	expectedFragment := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	if diff := cmp.Diff(record.Fragment, expectedFragment); diff != "" {
		t.Errorf("Fragment mismatch (-got +want):\n%s", diff)
	}
}

func TestPeekTLSRecord_TooLarge(t *testing.T) {
	// Record with length > 16384 (maximum TLS record size)
	tooLargeRecord := []byte{
		0x17,       // application_data
		0x03, 0x03, // TLS 1.2
		0x50, 0x00, // length: 20480 (> 16384)
	}

	peeker := &testPeeker{data: tooLargeRecord}
	_, _, _, err := PeekTLSRecord(peeker, nil)

	if err == nil {
		t.Error("Expected error for record that is too large")
	}
}

func TestPeekTLSRecord_Incomplete(t *testing.T) {
	// Incomplete header
	incompleteHeader := []byte{0x16, 0x03, 0x03}

	peeker := &testPeeker{data: incompleteHeader}
	_, tryAgain, processed, _ := PeekTLSRecord(peeker, nil)

	if processed != 0 {
		t.Errorf("Expected processed to be 0, got %d", processed)
	}

	if tryAgain != 5 {
		t.Errorf("Expected tryAgain to be 5, got %d", tryAgain)
	}
}

func TestDuplicateRecord(t *testing.T) {
	original := tlsmirror.TLSRecord{
		RecordType:            0x17,
		LegacyProtocolVersion: [2]byte{0x03, 0x03},
		RecordLength:          5,
		Fragment:              []byte{0x01, 0x02, 0x03, 0x04, 0x05},
	}

	duplicate := DuplicateRecord(original)

	// Check that the duplicate has the same values
	if diff := cmp.Diff(duplicate.RecordType, original.RecordType); diff != "" {
		t.Errorf("RecordType mismatch (-got +want):\n%s", diff)
	}

	if diff := cmp.Diff(duplicate.Fragment, original.Fragment); diff != "" {
		t.Errorf("Fragment mismatch (-got +want):\n%s", diff)
	}

	// Check that modifying duplicate doesn't affect original
	duplicate.Fragment[0] = 0xFF
	if original.Fragment[0] == 0xFF {
		t.Error("Modifying duplicate affected original fragment")
	}
}

func TestTLSRecordStreamReader(t *testing.T) {
	// Create multiple TLS records
	record1 := tlsmirror.TLSRecord{
		RecordType:            0x16,
		LegacyProtocolVersion: [2]byte{0x03, 0x03},
		RecordLength:          3,
		Fragment:              []byte{0x01, 0x02, 0x03},
	}

	record2 := tlsmirror.TLSRecord{
		RecordType:            0x17,
		LegacyProtocolVersion: [2]byte{0x03, 0x03},
		RecordLength:          4,
		Fragment:              []byte{0x0A, 0x0B, 0x0C, 0x0D},
	}

	// Pack both records into a stream
	var buf bytes.Buffer
	buf.Write(PackTLSRecord(record1))
	buf.Write(PackTLSRecord(record2))

	reader := bufio.NewReader(&buf)
	streamReader := NewTLSRecordStreamReader(reader)

	// Read first record
	read1, err := streamReader.ReadNextRecord()
	if err != nil {
		t.Fatalf("Failed to read first record: %v", err)
	}

	if read1.RecordType != record1.RecordType {
		t.Errorf("First record type mismatch: expected 0x%02x, got 0x%02x", record1.RecordType, read1.RecordType)
	}

	if diff := cmp.Diff(read1.Fragment, record1.Fragment); diff != "" {
		t.Errorf("First record fragment mismatch (-got +want):\n%s", diff)
	}

	// Read second record
	read2, err := streamReader.ReadNextRecord()
	if err != nil {
		t.Fatalf("Failed to read second record: %v", err)
	}

	if read2.RecordType != record2.RecordType {
		t.Errorf("Second record type mismatch: expected 0x%02x, got 0x%02x", record2.RecordType, read2.RecordType)
	}

	if diff := cmp.Diff(read2.Fragment, record2.Fragment); diff != "" {
		t.Errorf("Second record fragment mismatch (-got +want):\n%s", diff)
	}

	// Check consumed size
	expectedSize := int64(8 + 9) // (5 + 3) + (5 + 4)
	if streamReader.GetConsumedSize() != expectedSize {
		t.Errorf("Expected consumed size %d, got %d", expectedSize, streamReader.GetConsumedSize())
	}
}

func TestTLSRecordStreamWriter(t *testing.T) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	streamWriter := NewTLSRecordStreamWriter(writer)

	record := &tlsmirror.TLSRecord{
		RecordType:            0x17,
		LegacyProtocolVersion: [2]byte{0x03, 0x03},
		RecordLength:          5,
		Fragment:              []byte{0x01, 0x02, 0x03, 0x04, 0x05},
	}

	err := streamWriter.WriteRecord(record, false)
	if err != nil {
		t.Fatalf("Failed to write record: %v", err)
	}

	expected := []byte{
		0x17,       // record type
		0x03, 0x03, // protocol version
		0x00, 0x05, // length
		0x01, 0x02, 0x03, 0x04, 0x05, // fragment
	}

	if diff := cmp.Diff(buf.Bytes(), expected); diff != "" {
		t.Errorf("Written data mismatch (-got +want):\n%s", diff)
	}
}

func TestTLSRecordStreamRoundTrip(t *testing.T) {
	// Create test records
	records := []*tlsmirror.TLSRecord{
		{
			RecordType:            0x16,
			LegacyProtocolVersion: [2]byte{0x03, 0x03},
			RecordLength:          10,
			Fragment:              []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A},
		},
		{
			RecordType:            0x17,
			LegacyProtocolVersion: [2]byte{0x03, 0x04},
			RecordLength:          5,
			Fragment:              []byte{0x0A, 0x0B, 0x0C, 0x0D, 0x0E},
		},
	}

	// Write records
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	streamWriter := NewTLSRecordStreamWriter(writer)

	for _, record := range records {
		err := streamWriter.WriteRecord(record, false)
		if err != nil {
			t.Fatalf("Failed to write record: %v", err)
		}
	}

	// Read records back
	reader := bufio.NewReader(&buf)
	streamReader := NewTLSRecordStreamReader(reader)

	for i, expected := range records {
		read, err := streamReader.ReadNextRecord()
		if err != nil {
			t.Fatalf("Failed to read record %d: %v", i, err)
		}

		if read.RecordType != expected.RecordType {
			t.Errorf("Record %d: type mismatch: expected 0x%02x, got 0x%02x", i, expected.RecordType, read.RecordType)
		}

		if diff := cmp.Diff(read.Fragment, expected.Fragment); diff != "" {
			t.Errorf("Record %d: fragment mismatch (-got +want):\n%s", i, diff)
		}

		if diff := cmp.Diff(read.LegacyProtocolVersion[:], expected.LegacyProtocolVersion[:]); diff != "" {
			t.Errorf("Record %d: protocol version mismatch (-got +want):\n%s", i, diff)
		}
	}
}
