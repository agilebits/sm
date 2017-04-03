package secrets

import (
	"bytes"
	"testing"
)

func TestCanGenerateKey(t *testing.T) {
	svc := NewDevKeyService()

	const kid = "key1"
	key1, err := svc.GenerateKey(kid)
	if err != nil {
		t.Fatal("failed to generate key1:", err)
	}

	key2, err := svc.GenerateKey(kid)
	if err != nil {
		t.Fatal("failed to generate key2:", err)
	}

	if key1.KID != kid {
		t.Errorf("expected key1.KID = %q, got %q", kid, key1.KID)
	}

	if key2.KID != kid {
		t.Errorf("expected key2.KID = %q, got %q", kid, key2.KID)
	}

	if bytes.Compare(key1.RawKey, key2.RawKey) != 0 {
		t.Errorf("expect the same RawKey in key1 and key2")
	}

	key3, err := svc.GenerateKey("key3")
	if err != nil {
		t.Fatal("failed to generate key3:", err)
	}

	if bytes.Compare(key1.RawKey, key3.RawKey) == 0 {
		t.Errorf("expect the different RawKey in key1 and key3")
	}
}

func TestCanEncryptAndDecrypt(t *testing.T) {
	svc := NewDevKeyService()

	key, err := svc.GenerateKey("somekey")
	if err != nil {
		t.Fatal("failed to generate key:", err)
	}

	source := "Hello, World!"
	ciphertext, err := key.Encrypt([]byte(source))
	if err != nil {
		t.Fatal("failed to Encrypt:", err)
	}

	plaintext, err := key.Decrypt(ciphertext)
	if err != nil {
		t.Fatal("failed to Decrypt:", err)
	}

	if string(plaintext) != source {
		t.Errorf("encrypt/decrypt failed, expected %q but got %q", source, plaintext)
	}
}
