// +build !confonly

package command

//go:generate go run github.com/v2fly/v2ray-core/v4/common/errors/errorgen

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/v2fly/v2ray-core/v4/features"

	"google.golang.org/grpc"

	core "github.com/v2fly/v2ray-core/v4"
	"github.com/v2fly/v2ray-core/v4/app/observatory"
	"github.com/v2fly/v2ray-core/v4/common"
	"github.com/v2fly/v2ray-core/v4/features/extension"
)

type service struct {
	UnimplementedObservatoryServiceServer
	v *core.Instance

	observatory extension.Observatory
}

func (s *service) GetOutboundStatus(ctx context.Context, request *GetOutboundStatusRequest) (*GetOutboundStatusResponse, error) {
	var result proto.Message
	if request.Tag == "" {
		observeResult, err := s.observatory.GetObservation(ctx)
		if err != nil {
			newError("cannot get observation").Base(err)
		}
		result = observeResult
	} else {
		observeResult, err := common.Must2(s.observatory.(features.TaggedFeatures).GetFeaturesByTag(request.Tag)).(extension.Observatory).GetObservation(s.ctx)
		if err != nil {
			newError("cannot get observation").Base(err)
		}
		result = observeResult
	}
	retdata := result.(*observatory.ObservationResult)
	return &GetOutboundStatusResponse{
		Status: retdata,
	}, nil
}

func (s *service) Register(server *grpc.Server) {
	RegisterObservatoryServiceServer(server, s)
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, cfg interface{}) (interface{}, error) {
		s := core.MustFromContext(ctx)
		sv := &service{v: s}
		err := s.RequireFeatures(func(Observatory extension.Observatory) {
			sv.observatory = Observatory
		})
		if err != nil {
			return nil, err
		}
		return sv, nil
	}))
}
