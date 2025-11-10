package mirrorcrypto

import (
	"bytes"
	"testing"
)

func TestDeriveEncryptionKey(t *testing.T) {
	primaryKey := make([]byte, 32)
	clientRandom := make([]byte, 32)
	serverRandom := make([]byte, 32)

	// Fill with test data
	for i := 0; i < 32; i++ {
		primaryKey[i] = byte(i)
		clientRandom[i] = byte(i + 32)
		serverRandom[i] = byte(i + 64)
	}

	encKey, nonceMask, err := DeriveEncryptionKey(primaryKey, clientRandom, serverRandom, "test")
	if err != nil {
		t.Fatalf("DeriveEncryptionKey failed: %v", err)
	}

	// Verify key sizes
	if len(encKey) != 16 {
		t.Errorf("Expected encryption key length 16, got %d", len(encKey))
	}

	if len(nonceMask) != 12 {
		t.Errorf("Expected nonce mask length 12, got %d", len(nonceMask))
	}

	// Verify deterministic behavior - same inputs should produce same outputs
	encKey2, nonceMask2, err := DeriveEncryptionKey(primaryKey, clientRandom, serverRandom, "test")
	if err != nil {
		t.Fatalf("Second DeriveEncryptionKey failed: %v", err)
	}

	if !bytes.Equal(encKey, encKey2) {
		t.Error("Encryption keys are not deterministic")
	}

	if !bytes.Equal(nonceMask, nonceMask2) {
		t.Error("Nonce masks are not deterministic")
	}

	// Verify different tags produce different keys
	encKey3, nonceMask3, err := DeriveEncryptionKey(primaryKey, clientRandom, serverRandom, "different")
	if err != nil {
		t.Fatalf("Third DeriveEncryptionKey failed: %v", err)
	}

	if bytes.Equal(encKey, encKey3) {
		t.Error("Different tags should produce different encryption keys")
	}

	if bytes.Equal(nonceMask, nonceMask3) {
		t.Error("Different tags should produce different nonce masks")
	}
}

func TestDeriveEncryptionKey_InvalidPrimaryKeySize(t *testing.T) {
	invalidPrimaryKey := make([]byte, 16) // Wrong size
	clientRandom := make([]byte, 32)
	serverRandom := make([]byte, 32)

	_, _, err := DeriveEncryptionKey(invalidPrimaryKey, clientRandom, serverRandom, "test")
	if err == nil {
		t.Error("Expected error for invalid primary key size")
	}
}

func TestDeriveEncryptionKey_InvalidClientRandomSize(t *testing.T) {
	primaryKey := make([]byte, 32)
	invalidClientRandom := make([]byte, 16) // Wrong size
	serverRandom := make([]byte, 32)

	_, _, err := DeriveEncryptionKey(primaryKey, invalidClientRandom, serverRandom, "test")
	if err == nil {
		t.Error("Expected error for invalid client random size")
	}
}

func TestDeriveEncryptionKey_InvalidServerRandomSize(t *testing.T) {
	primaryKey := make([]byte, 32)
	clientRandom := make([]byte, 32)
	invalidServerRandom := make([]byte, 16) // Wrong size

	_, _, err := DeriveEncryptionKey(primaryKey, clientRandom, invalidServerRandom, "test")
	if err == nil {
		t.Error("Expected error for invalid server random size")
	}
}

func TestDeriveSecondaryKey(t *testing.T) {
	primaryKey := make([]byte, 32)
	for i := 0; i < 32; i++ {
		primaryKey[i] = byte(i)
	}

	secondaryKey, err := DeriveSecondaryKey(primaryKey, "test")
	if err != nil {
		t.Fatalf("DeriveSecondaryKey failed: %v", err)
	}

	// Verify key size
	if len(secondaryKey) != 16 {
		t.Errorf("Expected secondary key length 16, got %d", len(secondaryKey))
	}

	// Verify deterministic behavior
	secondaryKey2, err := DeriveSecondaryKey(primaryKey, "test")
	if err != nil {
		t.Fatalf("Second DeriveSecondaryKey failed: %v", err)
	}

	if !bytes.Equal(secondaryKey, secondaryKey2) {
		t.Error("Secondary keys are not deterministic")
	}

	// Verify different tags produce different keys
	secondaryKey3, err := DeriveSecondaryKey(primaryKey, "different")
	if err != nil {
		t.Fatalf("Third DeriveSecondaryKey failed: %v", err)
	}

	if bytes.Equal(secondaryKey, secondaryKey3) {
		t.Error("Different tags should produce different secondary keys")
	}
}

func TestDeriveSecondaryKey_InvalidPrimaryKeySize(t *testing.T) {
	invalidPrimaryKey := make([]byte, 16) // Wrong size

	_, err := DeriveSecondaryKey(invalidPrimaryKey, "test")
	if err == nil {
		t.Error("Expected error for invalid primary key size")
	}
}

func TestDeriveSequenceWatermarkingKey(t *testing.T) {
	primaryKey := make([]byte, 32)
	clientRandom := make([]byte, 32)
	serverRandom := make([]byte, 32)

	// Fill with test data
	for i := 0; i < 32; i++ {
		primaryKey[i] = byte(i)
		clientRandom[i] = byte(i + 32)
		serverRandom[i] = byte(i + 64)
	}

	encKey, nonceMask, err := DeriveSequenceWatermarkingKey(primaryKey, clientRandom, serverRandom, "test")
	if err != nil {
		t.Fatalf("DeriveSequenceWatermarkingKey failed: %v", err)
	}

	// Verify key sizes (different from DeriveEncryptionKey)
	if len(encKey) != 32 {
		t.Errorf("Expected watermarking encryption key length 32, got %d", len(encKey))
	}

	if len(nonceMask) != 24 {
		t.Errorf("Expected watermarking nonce mask length 24, got %d", len(nonceMask))
	}

	// Verify deterministic behavior
	encKey2, nonceMask2, err := DeriveSequenceWatermarkingKey(primaryKey, clientRandom, serverRandom, "test")
	if err != nil {
		t.Fatalf("Second DeriveSequenceWatermarkingKey failed: %v", err)
	}

	if !bytes.Equal(encKey, encKey2) {
		t.Error("Watermarking encryption keys are not deterministic")
	}

	if !bytes.Equal(nonceMask, nonceMask2) {
		t.Error("Watermarking nonce masks are not deterministic")
	}

	// Verify different tags produce different keys
	encKey3, nonceMask3, err := DeriveSequenceWatermarkingKey(primaryKey, clientRandom, serverRandom, "different")
	if err != nil {
		t.Fatalf("Third DeriveSequenceWatermarkingKey failed: %v", err)
	}

	if bytes.Equal(encKey, encKey3) {
		t.Error("Different tags should produce different watermarking encryption keys")
	}

	if bytes.Equal(nonceMask, nonceMask3) {
		t.Error("Different tags should produce different watermarking nonce masks")
	}

	// Verify watermarking keys are different from encryption keys
	regularEncKey, regularNonceMask, _ := DeriveEncryptionKey(primaryKey, clientRandom, serverRandom, "test")

	if bytes.Equal(encKey[:16], regularEncKey) {
		t.Error("Watermarking encryption key should be different from regular encryption key")
	}

	if bytes.Equal(nonceMask[:12], regularNonceMask) {
		t.Error("Watermarking nonce mask should be different from regular nonce mask")
	}
}

func TestDeriveSequenceWatermarkingKey_InvalidInputs(t *testing.T) {
	validKey := make([]byte, 32)
	validRandom := make([]byte, 32)
	invalidKey := make([]byte, 16)
	invalidRandom := make([]byte, 16)

	// Test invalid primary key
	_, _, err := DeriveSequenceWatermarkingKey(invalidKey, validRandom, validRandom, "test")
	if err == nil {
		t.Error("Expected error for invalid primary key size")
	}

	// Test invalid client random
	_, _, err = DeriveSequenceWatermarkingKey(validKey, invalidRandom, validRandom, "test")
	if err == nil {
		t.Error("Expected error for invalid client random size")
	}

	// Test invalid server random
	_, _, err = DeriveSequenceWatermarkingKey(validKey, validRandom, invalidRandom, "test")
	if err == nil {
		t.Error("Expected error for invalid server random size")
	}
}

func TestKeyDerivationUniqueness(t *testing.T) {
	// Verify that different inputs produce different outputs
	primaryKey1 := make([]byte, 32)
	primaryKey2 := make([]byte, 32)
	clientRandom := make([]byte, 32)
	serverRandom := make([]byte, 32)

	for i := 0; i < 32; i++ {
		primaryKey1[i] = byte(i)
		primaryKey2[i] = byte(i + 1) // Different primary key
		clientRandom[i] = byte(i + 32)
		serverRandom[i] = byte(i + 64)
	}

	encKey1, _, err := DeriveEncryptionKey(primaryKey1, clientRandom, serverRandom, "test")
	if err != nil {
		t.Fatalf("DeriveEncryptionKey failed: %v", err)
	}

	encKey2, _, err := DeriveEncryptionKey(primaryKey2, clientRandom, serverRandom, "test")
	if err != nil {
		t.Fatalf("DeriveEncryptionKey failed: %v", err)
	}

	if bytes.Equal(encKey1, encKey2) {
		t.Error("Different primary keys should produce different encryption keys")
	}
}
