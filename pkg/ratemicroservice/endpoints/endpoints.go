package endpoints

import (
	"context"
	"errors"
	"os"

	rate "github.com/A-Danylevych/btc-api/pkg/ratemicroservice"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

type Set struct {
	GetEndpoint           endpoint.Endpoint
	ServiceStatusEndpoint endpoint.Endpoint
}

func NewEndpointSet(svc rate.Service) Set {
	return Set{
		GetEndpoint:           MakeGetRateEndpoint(svc),
		ServiceStatusEndpoint: MakeServiceStatusEndpoint(svc),
	}
}

func MakeGetRateEndpoint(svc rate.Service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		rate, err := svc.GetRate(ctx)
		if err != nil {
			return GetResponse{Rate: rate, Err: err.Error()}, nil
		}
		return GetResponse{rate, ""}, nil
	}
}

func MakeServiceStatusEndpoint(svc rate.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(ServiceStatusRequest)
		code, err := svc.ServiceStatus(ctx)
		if err != nil {
			return ServiceStatusResponse{Code: code, Err: err.Error()}, nil
		}
		return ServiceStatusResponse{Code: code, Err: ""}, nil
	}
}

func (s *Set) Get(ctx context.Context) (float64, error) {
	resp, err := s.GetEndpoint(ctx, GetRequest{})
	if err != nil {
		return 0, err
	}
	getResp := resp.(GetResponse)
	if getResp.Err != "" {
		return 0, errors.New(getResp.Err)
	}
	return getResp.Rate, nil
}

func (s *Set) ServiceStatus(ctx context.Context) (int, error) {
	resp, err := s.ServiceStatusEndpoint(ctx, ServiceStatusRequest{})
	svcStatusResp := resp.(ServiceStatusResponse)
	if err != nil {
		return svcStatusResp.Code, err
	}
	if svcStatusResp.Err != "" {
		return svcStatusResp.Code, errors.New(svcStatusResp.Err)
	}
	return svcStatusResp.Code, nil
}

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}
