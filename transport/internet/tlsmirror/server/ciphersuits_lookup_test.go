package server

import (
	"testing"
)

func TestNewEmptyCipherSuiteLookuper(t *testing.T) {
	lookuper := newEmptyCipherSuiteLookuper()

	if lookuper == nil {
		t.Fatal("Expected non-nil lookuper")
	}

	// Empty lookuper should return false for any ciphersuite
	testCipherSuites := []uint16{0x0000, 0x1234, 0xFFFF, 0xC02F}
	for _, cs := range testCipherSuites {
		if lookuper.Lookup(cs) {
			t.Errorf("Empty lookuper should return false for ciphersuite 0x%04X", cs)
		}
	}
}

func TestNewCipherSuiteLookuperFromUint16Array(t *testing.T) {
	ciphersuites := []uint16{0xC02F, 0xC030, 0x1301, 0x1302}
	lookuper := newCipherSuiteLookuperFromUint16Array(ciphersuites)

	if lookuper == nil {
		t.Fatal("Expected non-nil lookuper")
	}

	// Test that all ciphersuites in the array are found
	for _, cs := range ciphersuites {
		if !lookuper.Lookup(cs) {
			t.Errorf("Expected lookuper to find ciphersuite 0x%04X", cs)
		}
	}

	// Test that ciphersuites not in the array are not found
	notIncluded := []uint16{0x0000, 0xFFFF, 0x1234}
	for _, cs := range notIncluded {
		if lookuper.Lookup(cs) {
			t.Errorf("Expected lookuper to not find ciphersuite 0x%04X", cs)
		}
	}
}

func TestNewCipherSuiteLookuperFromUint32Array(t *testing.T) {
	ciphersuites := []uint32{0xC02F, 0xC030, 0x1301, 0x1302}
	lookuper, err := newCipherSuiteLookuperFromUint32Array(ciphersuites)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if lookuper == nil {
		t.Fatal("Expected non-nil lookuper")
	}

	// Test that all ciphersuites are found
	for _, cs := range ciphersuites {
		if !lookuper.Lookup(uint16(cs)) {
			t.Errorf("Expected lookuper to find ciphersuite 0x%04X", cs)
		}
	}

	// Test that other ciphersuites are not found
	if lookuper.Lookup(0xFFFF) {
		t.Error("Expected lookuper to not find ciphersuite 0xFFFF")
	}
}

func TestNewCipherSuiteLookuperFromUint32Array_Empty(t *testing.T) {
	emptyArray := []uint32{}
	lookuper, err := newCipherSuiteLookuperFromUint32Array(emptyArray)

	if err == nil {
		t.Error("Expected error for empty ciphersuite array")
	}

	if lookuper != nil {
		t.Error("Expected nil lookuper for empty array")
	}
}

func TestNewCipherSuiteLookuperFromUint32Array_OutOfRange(t *testing.T) {
	// Value > 0xFFFF
	outOfRangeArray := []uint32{0xC02F, 0x10000}
	lookuper, err := newCipherSuiteLookuperFromUint32Array(outOfRangeArray)

	if err == nil {
		t.Error("Expected error for out-of-range ciphersuite value")
	}

	if lookuper != nil {
		t.Error("Expected nil lookuper for out-of-range value")
	}
}

func TestCipherSuiteLookup(t *testing.T) {
	// Common TLS 1.2 and TLS 1.3 ciphersuites
	ciphersuites := []uint16{
		0xC02F, // TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
		0xC030, // TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
		0x1301, // TLS_AES_128_GCM_SHA256
		0x1302, // TLS_AES_256_GCM_SHA384
		0x1303, // TLS_CHACHA20_POLY1305_SHA256
	}

	lookuper := newCipherSuiteLookuperFromUint16Array(ciphersuites)

	testCases := []struct {
		name        string
		ciphersuite uint16
		expected    bool
	}{
		{"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256", 0xC02F, true},
		{"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384", 0xC030, true},
		{"TLS_AES_128_GCM_SHA256", 0x1301, true},
		{"TLS_AES_256_GCM_SHA384", 0x1302, true},
		{"TLS_CHACHA20_POLY1305_SHA256", 0x1303, true},
		{"Not included - NULL", 0x0000, false},
		{"Not included - Other", 0x1234, false},
		{"Not included - MAX", 0xFFFF, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := lookuper.Lookup(tc.ciphersuite)
			if result != tc.expected {
				t.Errorf("Lookup(0x%04X) = %v, expected %v", tc.ciphersuite, result, tc.expected)
			}
		})
	}
}

func TestCipherSuiteLookuperDuplicates(t *testing.T) {
	// Test that duplicate entries are handled correctly
	ciphersuites := []uint16{0xC02F, 0xC02F, 0xC030, 0xC030}
	lookuper := newCipherSuiteLookuperFromUint16Array(ciphersuites)

	if !lookuper.Lookup(0xC02F) {
		t.Error("Expected to find 0xC02F")
	}

	if !lookuper.Lookup(0xC030) {
		t.Error("Expected to find 0xC030")
	}

	// Should still return false for non-included ciphersuites
	if lookuper.Lookup(0x1234) {
		t.Error("Expected not to find 0x1234")
	}
}

func TestCipherSuiteLookuperBoundaryValues(t *testing.T) {
	// Test boundary values
	ciphersuites := []uint16{0x0000, 0xFFFF}
	lookuper := newCipherSuiteLookuperFromUint16Array(ciphersuites)

	if !lookuper.Lookup(0x0000) {
		t.Error("Expected to find 0x0000")
	}

	if !lookuper.Lookup(0xFFFF) {
		t.Error("Expected to find 0xFFFF")
	}

	// Middle value should not be found
	if lookuper.Lookup(0x8000) {
		t.Error("Expected not to find 0x8000")
	}
}

func TestCipherSuiteLookuperSingleEntry(t *testing.T) {
	ciphersuites := []uint16{0x1301}
	lookuper := newCipherSuiteLookuperFromUint16Array(ciphersuites)

	if !lookuper.Lookup(0x1301) {
		t.Error("Expected to find 0x1301")
	}

	if lookuper.Lookup(0x1302) {
		t.Error("Expected not to find 0x1302")
	}
}

func TestUint32ToUint16Conversion(t *testing.T) {
	// Test valid conversions
	validCiphersuites := []uint32{0x0000, 0x1234, 0xFFFF}
	lookuper, err := newCipherSuiteLookuperFromUint32Array(validCiphersuites)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	for _, cs := range validCiphersuites {
		if !lookuper.Lookup(uint16(cs)) {
			t.Errorf("Expected to find ciphersuite 0x%04X", cs)
		}
	}
}

func TestUint32ToUint16ConversionMaxBoundary(t *testing.T) {
	// Test max valid value
	maxValid := []uint32{0xFFFF}
	lookuper, err := newCipherSuiteLookuperFromUint32Array(maxValid)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !lookuper.Lookup(0xFFFF) {
		t.Error("Expected to find 0xFFFF")
	}

	// Test just above max valid value
	justAboveMax := []uint32{0x10000}
	_, err = newCipherSuiteLookuperFromUint32Array(justAboveMax)

	if err == nil {
		t.Error("Expected error for value 0x10000")
	}
}
