package conf

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/v2fly/v2ray-core/v4/app/observatory/burst"
	"github.com/v2fly/v2ray-core/v4/app/observatory/multiObservatory"
	"github.com/v2fly/v2ray-core/v4/common/serial"

	"github.com/v2fly/v2ray-core/v4/app/observatory"
)

type ObservatoryConfig struct {
	SubjectSelector []string `json:"subjectSelector"`
}

func (o *ObservatoryConfig) Build() (proto.Message, error) {
	return &observatory.Config{SubjectSelector: o.SubjectSelector}, nil
}

type BurstObservatoryConfig struct {
	SubjectSelector []string `json:"subjectSelector"`
	// health check settings
	HealthCheck *healthCheckSettings `json:"pingConfig,omitempty"`
}

func (b BurstObservatoryConfig) Build() (proto.Message, error) {
	if result, err := b.HealthCheck.Build(); err != nil {
		return &burst.Config{SubjectSelector: b.SubjectSelector, PingConfig: result.(*burst.HealthPingConfig)}, nil
	} else {
		return nil, err
	}
}

type MultiObservatoryItem struct {
	MemberType string          `json:"type"`
	Tag        string          `json:"tag"`
	Value      json.RawMessage `json:"settings"`
}

type MultiObservatoryConfig struct {
	Observers []MultiObservatoryItem `json:"observers"`
}

func (o *MultiObservatoryConfig) Build() (proto.Message, error) {
	ret := &multiObservatory.Config{}
	for _, v := range o.Observers {
		switch v.MemberType {
		case "burst":
			var burstObservatoryConfig BurstObservatoryConfig
			err := json.Unmarshal(v.Value, &burstObservatoryConfig)
			if err != nil {
				return nil, err
			}
			burstObservatoryConfigPb, err := burstObservatoryConfig.Build()
			if err != nil {
				return nil, err
			}
			ret.Holders.Features[v.Tag] = serial.ToTypedMessage(burstObservatoryConfigPb)
			break
		case "default":
			fallthrough
		default:
			var observatoryConfig ObservatoryConfig
			err := json.Unmarshal(v.Value, &observatoryConfig)
			if err != nil {
				return nil, err
			}
			observatoryConfigPb, err := observatoryConfig.Build()
			if err != nil {
				return nil, err
			}
			ret.Holders.Features[v.Tag] = serial.ToTypedMessage(observatoryConfigPb)
		}
	}
	return ret, nil
}
