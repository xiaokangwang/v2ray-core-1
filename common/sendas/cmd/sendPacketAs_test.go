package sendascmd

import (
	"bytes"
	"net/netip"
	"strings"
	"testing"

	"github.com/v2fly/v2ray-core/v5/common/sendas"
)

func TestParseMethod(t *testing.T) {
	tests := []struct {
		input string
		want  sendas.Method
	}{
		{input: "raw", want: sendas.Method_RAW},
		{input: "RAW", want: sendas.Method_RAW},
		{input: "raw-ethernet", want: sendas.Method_RAW_ETHERNET},
		{input: "raw_ethernet", want: sendas.Method_RAW_ETHERNET},
	}

	for _, test := range tests {
		got, _, err := parseMethod(test.input)
		if err != nil {
			t.Fatalf("parseMethod(%q) failed: %v", test.input, err)
		}
		if got != test.want {
			t.Fatalf("parseMethod(%q) = %v, want %v", test.input, got, test.want)
		}
	}
}

func TestParseMethodRejectsUnknownValue(t *testing.T) {
	if _, _, err := parseMethod("bogus"); err == nil {
		t.Fatal("expected parseMethod to fail for unknown method")
	}
}

func TestDecodeHexPayload(t *testing.T) {
	payload, err := decodeHexPayload("0x68 65:6c-6c6f")
	if err != nil {
		t.Fatalf("decodeHexPayload failed: %v", err)
	}
	if string(payload) != "hello" {
		t.Fatalf("unexpected payload: %q", payload)
	}
}

func TestDecodeBase64Payload(t *testing.T) {
	payload, err := decodeBase64Payload("aGVsbG8=")
	if err != nil {
		t.Fatalf("decodeBase64Payload failed: %v", err)
	}
	if string(payload) != "hello" {
		t.Fatalf("unexpected payload: %q", payload)
	}
}

func TestParseCLIConfigValidatesRawEthernetRequirements(t *testing.T) {
	_, err := parseCLIConfig([]string{
		"-method", "raw-ethernet",
		"-src", "198.51.100.10:40000",
		"-dst", "203.0.113.20:53",
	}, bytes.NewBuffer(nil))
	if err == nil {
		t.Fatal("expected raw-ethernet config without interface and MAC to fail")
	}
	if !strings.Contains(err.Error(), "-interface") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseCLIConfigParsesGateway(t *testing.T) {
	cfg, err := parseCLIConfig([]string{
		"-src", "198.51.100.10:40000",
		"-dst", "203.0.113.20:53",
		"-gateway", "192.0.2.1",
	}, bytes.NewBuffer(nil))
	if err != nil {
		t.Fatalf("parseCLIConfig failed: %v", err)
	}
	if !cfg.gateway.IsValid() || cfg.gateway.String() != "192.0.2.1" {
		t.Fatalf("unexpected gateway: %v", cfg.gateway)
	}
}

func TestParseCLIConfigRejectsGatewayForRawEthernet(t *testing.T) {
	_, err := parseCLIConfig([]string{
		"-method", "raw-ethernet",
		"-src", "198.51.100.10:40000",
		"-dst", "203.0.113.20:53",
		"-gateway", "192.0.2.1",
		"-interface", "eth0",
		"-dst-mac", "02:00:00:00:00:02",
	}, bytes.NewBuffer(nil))
	if err == nil {
		t.Fatal("expected raw-ethernet gateway config to fail")
	}
	if !strings.Contains(err.Error(), "-gateway") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseCLIConfigRejectsMultiplePayloadFormats(t *testing.T) {
	_, err := parseCLIConfig([]string{
		"-src", "198.51.100.10:40000",
		"-dst", "203.0.113.20:53",
		"-payload", "hello",
		"-payload-base64", "aGVsbG8=",
	}, bytes.NewBuffer(nil))
	if err == nil {
		t.Fatal("expected multiple payload formats to fail")
	}
	if !strings.Contains(err.Error(), "only one") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPayloadForDefaultsToGeneratedProbe(t *testing.T) {
	cfg := cliConfig{
		src: netip.MustParseAddrPort("198.51.100.10:40000"),
		dst: netip.MustParseAddrPort("203.0.113.20:53"),
	}

	payload := cfg.payloadFor(2)
	if !strings.Contains(string(payload), "seq=3") {
		t.Fatalf("unexpected payload: %q", payload)
	}
	if !strings.Contains(string(payload), cfg.src.String()) {
		t.Fatalf("payload missing source: %q", payload)
	}
}
