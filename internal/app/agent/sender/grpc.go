package sender

import (
	"context"
	"github.com/vladislaoramos/alemetric/internal/entity"
	grpcTool "github.com/vladislaoramos/alemetric/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
)

type GRPCAgent struct {
	serverURL string
	//client    *resty.Client
	addr net.IP

	client grpcTool.MetricsToolClient
}

func NewGRPCAgent(url string, ip string) *GRPCAgent {
	return &GRPCAgent{serverURL: url, addr: []byte(ip)}
}

func (g *GRPCAgent) Connect() (func(), error) {
	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	g.client = grpcTool.NewMetricsToolClient(conn)

	return func() {
		g.client = nil
		conn.Close()
	}, nil
}

func (g *GRPCAgent) SendMetrics(
	ctx context.Context,
	metricsName string,
	metricsType string,
	delta *entity.Counter,
	value *entity.Gauge,
) error {
	item := entity.Metrics{
		ID:    metricsName,
		MType: metricsType,
		Delta: delta,
		Value: value,
	}

	_, err := g.client.StoreMetrics(ctx, item.AsProto())
	if err != nil {
		return err
	}

	return nil
}

func (g *GRPCAgent) SendSeveralMetrics(_ context.Context, _ []entity.Metrics) error {
	return nil
}
