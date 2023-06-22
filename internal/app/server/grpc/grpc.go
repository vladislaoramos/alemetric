package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/vladislaoramos/alemetric/internal/entity"
	"github.com/vladislaoramos/alemetric/internal/repo"
	"github.com/vladislaoramos/alemetric/internal/usecase"
	logger "github.com/vladislaoramos/alemetric/pkg/log"
	grpcTool "github.com/vladislaoramos/alemetric/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"time"
)

type GRPCServer struct {
	grpcTool.UnimplementedMetricsToolServer

	Server  *grpc.Server
	log     logger.LogInterface
	storage usecase.MetricsRepo

	writeFileDuration       time.Duration
	writeToFileWithDuration bool
	syncWriteFile           bool
	asyncWriteFile          bool
	C                       chan struct{}

	checkDataSign bool
	encryptionKey string
}

func NewServer(
	repo usecase.MetricsRepo,
	l logger.LogInterface,
	options ...OptionFunc,
) *GRPCServer {
	g := grpc.NewServer()

	s := &GRPCServer{
		Server:  g,
		log:     l,
		storage: repo,
	}

	for _, o := range options {
		o(s)
	}

	if s.writeToFileWithDuration {
		go func() {
			ticker := time.NewTicker(s.writeFileDuration)
			for {
				<-ticker.C
				s.C <- struct{}{}
			}
		}()
	}

	if s.writeToFileWithDuration || s.asyncWriteFile {
		s.C = make(chan struct{}, 1)
		go s.saveStorage()
	}

	grpcTool.RegisterMetricsToolServer(g, s)

	return s
}

func (s *GRPCServer) saveStorage() {
	for {
		<-s.C
		err := s.storage.StoreAll()
		if err != nil {
			s.log.Error(fmt.Sprintf("error while writing to storage: %s", err))
		} else {
			s.log.Info("successful saving of metrics")
		}
	}
}

func (s *GRPCServer) Start(port int) error {
	addr := net.TCPAddr{Port: port}
	listen, err := net.Listen("tcp", addr.String())
	if err != nil {
		return err
	}

	s.log.Info(fmt.Sprintf("Listening gRPC port %d", port))
	if err = s.Server.Serve(listen); err != nil {
		return err
	}

	return nil
}

func (s *GRPCServer) GetMetricsNames(
	ctx context.Context, _ *emptypb.Empty) (*grpcTool.MetricsNames, error) {
	names := s.storage.GetMetricsNames(ctx)

	res := grpcTool.MetricsNames{
		Response: &grpcTool.MetricsNames_Metrics{
			Metrics: &grpcTool.MetricsNamesList{Items: names},
		},
	}

	return &res, nil
}

func (s *GRPCServer) StoreMetrics(
	ctx context.Context, in *grpcTool.Metrics) (*emptypb.Empty, error) {
	metrics := entity.FromProto(in)
	if s.checkDataSign && !metrics.CheckDataSign(s.encryptionKey) {
		return &emptypb.Empty{}, fmt.Errorf("invalid metrics hash")
	}

	switch metrics.MType {
	case usecase.Gauge:
		if err := s.storage.StoreMetrics(ctx, metrics); err != nil {
			if errors.Is(err, repo.ErrNotFound) {
				return &emptypb.Empty{}, fmt.Errorf("store metrics: %w", repo.ErrNotFound)
			}
			return &emptypb.Empty{}, fmt.Errorf("error store metrics: %w", err)
		}
	case usecase.Counter:
		oldMetric, err := s.storage.GetMetrics(ctx, metrics.ID)
		if err != nil && !errors.Is(err, repo.ErrNotFound) {
			return &emptypb.Empty{}, fmt.Errorf("error getting metrics: %w", err)
		} else if errors.Is(err, repo.ErrNotFound) {
			var oldDelta entity.Counter
			newDelta := oldDelta + *metrics.Delta
			metrics.Delta = &newDelta
		} else {
			delta := *oldMetric.Delta + *metrics.Delta
			metrics.Delta = &delta
		}

		metrics.SignData("server", s.encryptionKey)

		if err = s.storage.StoreMetrics(ctx, metrics); err != nil {
			return &emptypb.Empty{}, fmt.Errorf("error storing metrics: %w", err)
		}

	default:
		return &emptypb.Empty{}, usecase.ErrNotImplemented
	}

	if s.asyncWriteFile {
		s.C <- struct{}{}
	}

	if s.syncWriteFile {
		err := s.storage.StoreAll()
		if err != nil {
			return &emptypb.Empty{}, fmt.Errorf("error storing all metrics: %w", err)
		}
	}

	return &emptypb.Empty{}, nil
}

func (s *GRPCServer) GetMetrics(
	ctx context.Context, in *grpcTool.Metrics) (*grpcTool.MetricsResponse, error,
) {
	m := entity.FromProto(in)
	res, err := s.storage.GetMetrics(ctx, m.ID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, usecase.ErrNotFound
		}
		return nil, fmt.Errorf("error getting metrics: %w", err)
	}

	if s.encryptionKey != "" && res.Hash == "" {
		res.SignData("server", s.encryptionKey)
	}

	resp := grpcTool.MetricsResponse{
		Response: &grpcTool.MetricsResponse_Metrics{
			Metrics: res.AsProto(),
		},
	}

	return &resp, nil
}
