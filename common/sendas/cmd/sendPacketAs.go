package sendascmd

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/netip"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/v2fly/v2ray-core/v5/common/sendas"
	"github.com/v2fly/v2ray-core/v5/common/sendas/raw"
	"github.com/v2fly/v2ray-core/v5/main/commands/all/engineering"
	"github.com/v2fly/v2ray-core/v5/main/commands/base"
)

type cliConfig struct {
	method          sendas.Method
	methodName      string
	interfaceName   string
	src             netip.AddrPort
	dst             netip.AddrPort
	gateway         netip.Addr
	srcMAC          []byte
	dstMAC          []byte
	explicitPayload []byte
	count           int
	interval        time.Duration
	quiet           bool
}

const sendPacketAsCommandName = "send-packet-as"

var cmdSendPacketAs = &base.Command{
	UsageLine:   "{{.Exec}} engineering send-packet-as [options]",
	Short:       "send UDP packets as a selected source address",
	CustomFlags: true,
	Long: `
Send UDP packets as the supplied source address through the sendas raw sender.

Raw mode sends an IP packet directly. With -gateway, route lookup uses that
next hop while the packet header keeps the original destination.
Raw-ethernet mode sends a full Ethernet frame and requires -interface and
-dst-mac.

This is an engineering validation tool for the sendas send-as path.
Linux raw sockets usually require root or CAP_NET_RAW.

Usage:
	{{.Exec}} engineering send-packet-as --src 198.51.100.10:40000 --dst 203.0.113.20:53 [options]

Examples:
	{{.Exec}} engineering send-packet-as --method raw --interface eth0 --src 198.51.100.10:40000 --dst 203.0.113.20:53 --gateway 192.0.2.1 --count 5 --interval 200ms
	{{.Exec}} engineering send-packet-as --method raw-ethernet --interface eth0 --src 198.51.100.10:40000 --dst 203.0.113.20:53 --dst-mac 02:00:00:00:00:02 --payload "sendas smoke test"
	{{.Exec}} engineering send-packet-as --method raw --src 198.51.100.10:40000 --dst 203.0.113.20:53 --payload-base64 c2VuZGFzIHNtb2tlIHRlc3Q=
`,
	Run: executeSendPacketAs,
}

func init() {
	engineering.AddCommand(cmdSendPacketAs)
}

func executeSendPacketAs(cmd *base.Command, args []string) {
	if err := run(args, os.Stdout, os.Stderr); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		base.Fatalf("send-packet-as: %v", err)
	}
}

func run(args []string, stdout io.Writer, stderr io.Writer) (err error) {
	cfg, err := parseCLIConfig(args, stderr)
	if err != nil {
		return err
	}

	if runtime.GOOS != "linux" {
		fmt.Fprintf(stderr, "warning: current platform is %s; raw sendas packet forging is implemented for linux\n", runtime.GOOS)
	}

	fmt.Fprintf(stderr, "sending %d packet(s) method=%s interface=%q src=%s dst=%s payload=%dB\n",
		cfg.count, cfg.methodName, cfg.interfaceName, cfg.src, cfg.dst, len(cfg.payloadFor(0)))
	if cfg.gateway.IsValid() {
		fmt.Fprintf(stderr, "gateway override: %s\n", cfg.gateway)
	}
	if len(cfg.srcMAC) != 0 {
		fmt.Fprintf(stderr, "source MAC override: %s\n", net.HardwareAddr(cfg.srcMAC))
	}
	if len(cfg.dstMAC) != 0 {
		fmt.Fprintf(stderr, "destination MAC override: %s\n", net.HardwareAddr(cfg.dstMAC))
	}

	sender := raw.NewUDPPacketSender(&sendas.Config{
		Interface: cfg.interfaceName,
		Method:    cfg.method,
	})
	defer func() {
		if closeErr := sender.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close sender: %w", closeErr)
		}
	}()

	for i := 0; i < cfg.count; i++ {
		payload := cfg.payloadFor(i)
		started := time.Now()
		if sendErr := sender.SendToWithGateway(cfg.src, cfg.dst, cfg.gateway, cfg.srcMAC, cfg.dstMAC, payload); sendErr != nil {
			return fmt.Errorf("send packet %d/%d: %w", i+1, cfg.count, sendErr)
		}

		if !cfg.quiet {
			fmt.Fprintf(stdout, "sent packet %d/%d in %s (%d bytes)\n", i+1, cfg.count, time.Since(started), len(payload))
		}

		if i+1 < cfg.count && cfg.interval > 0 {
			time.Sleep(cfg.interval)
		}
	}

	return nil
}

func parseCLIConfig(args []string, stderr io.Writer) (cliConfig, error) {
	var cfg cliConfig
	var method string
	var src string
	var dst string
	var gateway string
	var srcMAC string
	var dstMAC string
	var payloadText string
	var payloadHex string
	var payloadBase64 string

	fs := flag.NewFlagSet(sendPacketAsCommandName, flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		printUsage(fs, stderr)
	}

	fs.StringVar(&method, "method", "raw", "send method: raw or raw-ethernet")
	fs.StringVar(&cfg.interfaceName, "interface", "", "interface to bind or transmit on")
	fs.StringVar(&src, "src", "", "send as addr:port, e.g. 198.51.100.10:40000")
	fs.StringVar(&dst, "dst", "", "destination addr:port, e.g. 203.0.113.20:53")
	fs.StringVar(&gateway, "gateway", "", "route raw packets via this next-hop IP while preserving the original destination in the packet header")
	fs.StringVar(&srcMAC, "src-mac", "", "source MAC override for raw-ethernet mode")
	fs.StringVar(&dstMAC, "dst-mac", "", "destination MAC for raw-ethernet mode")
	fs.StringVar(&payloadText, "payload", "", "UTF-8 payload text; default is a generated probe line")
	fs.StringVar(&payloadHex, "payload-hex", "", "hex payload, separators like ':' '-' and spaces are allowed")
	fs.StringVar(&payloadBase64, "payload-base64", "", "base64 payload; standard and raw base64 encodings are accepted")
	fs.IntVar(&cfg.count, "count", 1, "number of packets to send")
	fs.DurationVar(&cfg.interval, "interval", 0, "delay between packets, e.g. 200ms")
	fs.BoolVar(&cfg.quiet, "quiet", false, "suppress per-packet success output")

	if err := fs.Parse(args); err != nil {
		return cliConfig{}, err
	}
	if len(fs.Args()) != 0 {
		return cliConfig{}, fmt.Errorf("unexpected positional arguments: %s", strings.Join(fs.Args(), " "))
	}

	parsedMethod, methodName, err := parseMethod(method)
	if err != nil {
		return cliConfig{}, err
	}
	cfg.method = parsedMethod
	cfg.methodName = methodName

	if cfg.count <= 0 {
		return cliConfig{}, fmt.Errorf("count must be greater than zero")
	}
	if cfg.interval < 0 {
		return cliConfig{}, fmt.Errorf("interval must not be negative")
	}
	if src == "" {
		return cliConfig{}, fmt.Errorf("missing required -src")
	}
	if dst == "" {
		return cliConfig{}, fmt.Errorf("missing required -dst")
	}
	payloadSources := 0
	if payloadText != "" {
		payloadSources++
	}
	if payloadHex != "" {
		payloadSources++
	}
	if payloadBase64 != "" {
		payloadSources++
	}
	if payloadSources > 1 {
		return cliConfig{}, fmt.Errorf("provide only one of -payload, -payload-hex, or -payload-base64")
	}

	cfg.src, err = netip.ParseAddrPort(src)
	if err != nil {
		return cliConfig{}, fmt.Errorf("parse -src: %w", err)
	}
	cfg.dst, err = netip.ParseAddrPort(dst)
	if err != nil {
		return cliConfig{}, fmt.Errorf("parse -dst: %w", err)
	}
	if gateway != "" {
		cfg.gateway, err = netip.ParseAddr(gateway)
		if err != nil {
			return cliConfig{}, fmt.Errorf("parse -gateway: %w", err)
		}
	}

	if srcMAC != "" {
		cfg.srcMAC, err = parseMACAddress(srcMAC)
		if err != nil {
			return cliConfig{}, fmt.Errorf("parse -src-mac: %w", err)
		}
	}
	if dstMAC != "" {
		cfg.dstMAC, err = parseMACAddress(dstMAC)
		if err != nil {
			return cliConfig{}, fmt.Errorf("parse -dst-mac: %w", err)
		}
	}
	if payloadHex != "" {
		cfg.explicitPayload, err = decodeHexPayload(payloadHex)
		if err != nil {
			return cliConfig{}, fmt.Errorf("parse -payload-hex: %w", err)
		}
	} else if payloadBase64 != "" {
		cfg.explicitPayload, err = decodeBase64Payload(payloadBase64)
		if err != nil {
			return cliConfig{}, fmt.Errorf("parse -payload-base64: %w", err)
		}
	} else if payloadText != "" {
		cfg.explicitPayload = []byte(payloadText)
	}

	if cfg.method == sendas.Method_RAW_ETHERNET {
		if cfg.gateway.IsValid() {
			return cliConfig{}, fmt.Errorf("-gateway is only supported in raw mode")
		}
		if cfg.interfaceName == "" {
			return cliConfig{}, fmt.Errorf("-interface is required for raw-ethernet mode")
		}
		if len(cfg.dstMAC) == 0 {
			return cliConfig{}, fmt.Errorf("-dst-mac is required for raw-ethernet mode")
		}
	}

	return cfg, nil
}

func printUsage(fs *flag.FlagSet, out io.Writer) {
	fmt.Fprintf(out, "Usage of %s:\n", fs.Name())
	fmt.Fprintf(out, "  %s engineering %s --src 198.51.100.10:40000 --dst 203.0.113.20:53 [options]\n\n", base.CommandEnv.Exec, sendPacketAsCommandName)
	fmt.Fprintln(out, "This engineering tool sends UDP packets as the supplied source address through the sendas raw sender.")
	fmt.Fprintln(out, "Raw mode sends an IP packet directly. With -gateway, route lookup uses that next hop while the packet header keeps the original destination.")
	fmt.Fprintln(out, "Raw-ethernet mode sends a full Ethernet frame and requires -interface and -dst-mac.")
	fmt.Fprintln(out, "Linux raw sockets usually require root or CAP_NET_RAW.")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Examples:")
	fmt.Fprintf(out, "  %s engineering %s --method raw --interface eth0 --src 198.51.100.10:40000 --dst 203.0.113.20:53 --gateway 192.0.2.1 --count 5 --interval 200ms\n", base.CommandEnv.Exec, sendPacketAsCommandName)
	fmt.Fprintf(out, "  %s engineering %s --method raw-ethernet --interface eth0 --src 198.51.100.10:40000 --dst 203.0.113.20:53 --dst-mac 02:00:00:00:00:02 --payload 'sendas smoke test'\n", base.CommandEnv.Exec, sendPacketAsCommandName)
	fmt.Fprintf(out, "  %s engineering %s --method raw --src 198.51.100.10:40000 --dst 203.0.113.20:53 --payload-base64 c2VuZGFzIHNtb2tlIHRlc3Q=\n", base.CommandEnv.Exec, sendPacketAsCommandName)
	fmt.Fprintln(out)
	fs.PrintDefaults()
}

func parseMethod(value string) (sendas.Method, string, error) {
	switch normalized := strings.ToLower(strings.TrimSpace(value)); normalized {
	case "raw":
		return sendas.Method_RAW, "raw", nil
	case "raw-ethernet", "raw_ethernet", "ethernet":
		return sendas.Method_RAW_ETHERNET, "raw-ethernet", nil
	default:
		return sendas.Method_FAIL, "", fmt.Errorf("unsupported method %q", value)
	}
}

func parseMACAddress(value string) ([]byte, error) {
	mac, err := net.ParseMAC(strings.TrimSpace(value))
	if err != nil {
		return nil, err
	}
	if len(mac) != 6 {
		return nil, fmt.Errorf("expected 6-byte MAC, got %d bytes", len(mac))
	}
	return []byte(mac), nil
}

func decodeHexPayload(value string) ([]byte, error) {
	cleaned := normalizeHexString(value)
	if cleaned == "" {
		return nil, fmt.Errorf("empty hex payload")
	}
	if len(cleaned)%2 != 0 {
		return nil, fmt.Errorf("hex payload must contain an even number of digits")
	}

	decoded, err := hex.DecodeString(cleaned)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

func decodeBase64Payload(value string) ([]byte, error) {
	cleaned := strings.TrimSpace(value)
	if cleaned == "" {
		return nil, fmt.Errorf("empty base64 payload")
	}

	decoded, err := base64.StdEncoding.DecodeString(cleaned)
	if err == nil {
		return decoded, nil
	}

	decoded, rawErr := base64.RawStdEncoding.DecodeString(cleaned)
	if rawErr == nil {
		return decoded, nil
	}

	return nil, err
}

func normalizeHexString(value string) string {
	cleaned := strings.TrimSpace(strings.ToLower(value))
	cleaned = strings.TrimPrefix(cleaned, "0x")
	replacer := strings.NewReplacer(":", "", "-", "", " ", "", "\t", "", "\n", "", "\r", "")
	return replacer.Replace(cleaned)
}

func (c cliConfig) payloadFor(seq int) []byte {
	if len(c.explicitPayload) != 0 {
		return c.explicitPayload
	}

	return []byte(fmt.Sprintf(
		"sendas probe seq=%d src=%s dst=%s time=%s",
		seq+1,
		c.src,
		c.dst,
		time.Now().UTC().Format(time.RFC3339Nano),
	))
}
