package popcount

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"testing"
)

func Test_Popcount(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "reasonable config",
			args: args{
				config: &Config{
					ExpansionPayloadLength: 16,
					Seed:                   "popcount",
					AddedLength:            16,
					DiversifierLength:      8,
				},
			},
		},
		{
			name: "so many zeros",
			args: args{
				config: &Config{
					ExpansionPayloadLength: 16,
					Seed:                   "popcount",
					AddedLength:            1024,
					DiversifierLength:      8,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < 100; i++ {
				inputLen := 7*i + 16
				t.Run(fmt.Sprintf("%v bytes input", inputLen), func(t *testing.T) {
					input := make([]byte, inputLen)
					io.ReadFull(rand.Reader, input)
					encodedData, err := popcountEncode(input, tt.args.config)
					if err != nil {
						t.Error(err)
					}
					if len(encodedData) != inputLen+int(tt.args.config.AddedLength+tt.args.config.DiversifierLength) {
						t.Error("encoded data length not match")
					}
					decodedData, err := popcountDecode(encodedData, tt.args.config)
					if err != nil {
						t.Error(err)
					}
					if len(decodedData) != inputLen {
						t.Error("decoded data length not match")
					}
					if !bytes.Equal(input, decodedData) {
						t.Error("decoded data not match")
					}
				})
			}

		})
	}
}
