package mirrorcommon

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUnpackTLSClientHello(t *testing.T) {
	// Minimal valid ClientHello structure
	// HandshakeType: 0x01 (ClientHello)
	// Length: 0x000026 (38 bytes)
	// Version: 0x0303 (TLS 1.2)
	// ClientRandom: 32 bytes
	clientHelloData := []byte{
		0x01,             // HandshakeType: ClientHello
		0x00, 0x00, 0x26, // Length: 38 bytes (not including handshake header)
		0x03, 0x03, // Version: TLS 1.2
		// ClientRandom: 32 bytes
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
	}

	clientHello, err := UnpackTLSClientHello(clientHelloData)
	if err != nil {
		t.Fatalf("Failed to unpack ClientHello: %v", err)
	}

	if clientHello.HandshakeType != 0x01 {
		t.Errorf("Expected HandshakeType 0x01, got 0x%02x", clientHello.HandshakeType)
	}

	if clientHello.Version != 0x0303 {
		t.Errorf("Expected Version 0x0303, got 0x%04x", clientHello.Version)
	}

	expectedRandom := [32]byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
	}

	if diff := cmp.Diff(clientHello.ClientRandom[:], expectedRandom[:]); diff != "" {
		t.Errorf("ClientRandom mismatch (-got +want):\n%s", diff)
	}
}

func TestUnpackTLSServerHello(t *testing.T) {
	// Minimal valid ServerHello structure
	// HandshakeType: 0x02 (ServerHello)
	// Length: 0x000027 (39 bytes)
	// Version: 0x0303 (TLS 1.2)
	// ServerRandom: 32 bytes
	// SessionIDLength: 0x00
	// CipherSuite: 0xC02F (TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256)
	serverHelloData := []byte{
		0x02,             // HandshakeType: ServerHello
		0x00, 0x00, 0x27, // Length: 39 bytes
		0x03, 0x03, // Version: TLS 1.2
		// ServerRandom: 32 bytes
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
		0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F, 0x30,
		0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,
		0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F, 0x40,
		0x00,       // SessionIDLength: 0
		0xC0, 0x2F, // CipherSuite: TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
	}

	serverHello, err := UnpackTLSServerHello(serverHelloData)
	if err != nil {
		t.Fatalf("Failed to unpack ServerHello: %v", err)
	}

	if serverHello.HandshakeType != 0x02 {
		t.Errorf("Expected HandshakeType 0x02, got 0x%02x", serverHello.HandshakeType)
	}

	if serverHello.Version != 0x0303 {
		t.Errorf("Expected Version 0x0303, got 0x%04x", serverHello.Version)
	}

	expectedRandom := [32]byte{
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
		0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F, 0x30,
		0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,
		0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F, 0x40,
	}

	if diff := cmp.Diff(serverHello.ServerRandom[:], expectedRandom[:]); diff != "" {
		t.Errorf("ServerRandom mismatch (-got +want):\n%s", diff)
	}

	if serverHello.SessionIDLength != 0 {
		t.Errorf("Expected SessionIDLength 0, got %d", serverHello.SessionIDLength)
	}

	if serverHello.CipherSuite != 0xC02F {
		t.Errorf("Expected CipherSuite 0xC02F, got 0x%04x", serverHello.CipherSuite)
	}
}

func TestUnpackTLSServerHello_WithSessionID(t *testing.T) {
	// ServerHello with session ID
	sessionID := []byte{0xAA, 0xBB, 0xCC, 0xDD}
	serverHelloData := []byte{
		0x02,             // HandshakeType: ServerHello
		0x00, 0x00, 0x2B, // Length: 43 bytes
		0x03, 0x03, // Version: TLS 1.2
		// ServerRandom: 32 bytes
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
		0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F, 0x30,
		0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,
		0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F, 0x40,
		0x04,                   // SessionIDLength: 4
		0xAA, 0xBB, 0xCC, 0xDD, // SessionID
		0xC0, 0x2F, // CipherSuite
	}

	serverHello, err := UnpackTLSServerHello(serverHelloData)
	if err != nil {
		t.Fatalf("Failed to unpack ServerHello: %v", err)
	}

	if serverHello.SessionIDLength != 4 {
		t.Errorf("Expected SessionIDLength 4, got %d", serverHello.SessionIDLength)
	}

	if diff := cmp.Diff(serverHello.SessionID, sessionID); diff != "" {
		t.Errorf("SessionID mismatch (-got +want):\n%s", diff)
	}

	if serverHello.CipherSuite != 0xC02F {
		t.Errorf("Expected CipherSuite 0xC02F, got 0x%04x", serverHello.CipherSuite)
	}
}

func TestUnpackTLSClientHello_Invalid(t *testing.T) {
	// Test with incomplete data
	incompleteData := []byte{0x01, 0x00, 0x00}

	_, err := UnpackTLSClientHello(incompleteData)
	if err == nil {
		t.Error("Expected error for incomplete ClientHello data")
	}
}

func TestUnpackTLSServerHello_Invalid(t *testing.T) {
	// Test with incomplete data
	incompleteData := []byte{0x02, 0x00, 0x00}

	_, err := UnpackTLSServerHello(incompleteData)
	if err == nil {
		t.Error("Expected error for incomplete ServerHello data")
	}
}
