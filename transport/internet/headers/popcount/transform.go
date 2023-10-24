package popcount

import (
	"crypto/rand"
	"hash/crc64"
	"io"
	mrand "math/rand"
)

const (
	bitKind_Primary = 1 + iota
	bitKind_Secondary
)

type placement struct {
	kind int
}

func generateDisplacementMap(randomSeed uint64, primaryInput, secondaryInput int) []placement {
	randSource := mrand.NewSource(int64(randomSeed))
	randDisplacementer := mrand.New(randSource)

	placementMap := make([]placement, (primaryInput+secondaryInput)*8)
	for i := range placementMap {
		if i < primaryInput*8 {
			placementMap[i] = placement{
				kind: bitKind_Primary,
			}
		} else {
			placementMap[i] = placement{
				kind: bitKind_Secondary,
			}
		}
	}
	randDisplacementer.Shuffle(len(placementMap), func(i, j int) {
		placementMap[i], placementMap[j] = placementMap[j], placementMap[i]
	})
	return placementMap
}

func bitwiseBlend(primaryInput, secondaryInput []byte, placementMap []placement) []byte {
	output := make([]byte, len(primaryInput)+len(secondaryInput))
	var primaryCurrent, secondaryCurrent int
	for i, p := range placementMap {
		if p.kind == bitKind_Primary {
			output[i/8] |= ((primaryInput[primaryCurrent/8] >> (primaryCurrent % 8)) & 0x01) << (i % 8)
			primaryCurrent++
		} else {
			output[i/8] |= ((secondaryInput[secondaryCurrent/8] >> (secondaryCurrent % 8)) & 0x01) << (i % 8)
			secondaryCurrent++
		}
	}
	return output
}

func bitwiseUnwindBlend(input []byte, primaryInputLen, secondaryInputLen int, placementMap []placement) ([]byte, []byte) {
	primaryOutput := make([]byte, primaryInputLen)
	secondaryOutput := make([]byte, secondaryInputLen)
	var primaryCurrent, secondaryCurrent int
	for i, p := range placementMap {
		if p.kind == bitKind_Primary {
			primaryOutput[primaryCurrent/8] |= ((input[i/8] >> (i % 8)) & 0x01) << (primaryCurrent % 8)
			primaryCurrent++
		} else {
			secondaryOutput[secondaryCurrent/8] |= ((input[i/8] >> (i % 8)) & 0x01) << (secondaryCurrent % 8)
			secondaryCurrent++
		}
	}
	return primaryOutput, secondaryOutput
}

func popcountEncode(input []byte, config *Config) ([]byte, error) {
	diversifier := make([]byte, config.DiversifierLength)
	{
		_, err := io.ReadFull(rand.Reader, diversifier)
		if err != nil {
			return nil, newError("failed to generate diversifier").Base(err)
		}
	}
	var popcountReducedInput []byte
	{
		randomSource := crc64.Checksum(diversifier, crc64.MakeTable(crc64.ECMA)) ^
			crc64.Checksum([]byte(config.Seed), crc64.MakeTable(crc64.ECMA))
		filament := make([]byte, config.AddedLength)
		placementMap := generateDisplacementMap(randomSource, len(input), len(filament))
		popcountReducedInput = bitwiseBlend(input, filament, placementMap)
	}
	var output []byte
	{
		randomSource := crc64.Checksum([]byte(config.Seed), crc64.MakeTable(crc64.ECMA))
		placementMap := generateDisplacementMap(randomSource, len(popcountReducedInput), len(diversifier))
		output = bitwiseBlend(popcountReducedInput, diversifier, placementMap)
	}
	return output, nil
}

func popcountDecode(input []byte, config *Config) ([]byte, error) {
	var popcountReducedInput, diversifier []byte
	{
		randomSource := crc64.Checksum([]byte(config.Seed), crc64.MakeTable(crc64.ECMA))
		placementMap := generateDisplacementMap(randomSource, len(input)-int(config.DiversifierLength), int(config.DiversifierLength))
		popcountReducedInput, diversifier =
			bitwiseUnwindBlend(input, len(input)-int(config.DiversifierLength), int(config.DiversifierLength), placementMap)
	}
	var output []byte
	{
		randomSource := crc64.Checksum(diversifier, crc64.MakeTable(crc64.ECMA)) ^
			crc64.Checksum([]byte(config.Seed), crc64.MakeTable(crc64.ECMA))
		filament := make([]byte, config.AddedLength)
		placementMap := generateDisplacementMap(randomSource, len(popcountReducedInput)-len(filament), len(filament))
		output, _ = bitwiseUnwindBlend(popcountReducedInput, len(popcountReducedInput)-len(filament), len(filament), placementMap)
	}
	return output, nil
}
