package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jdy879526487/cfg"
	proto "github.com/jdy879526487/cfg/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

// 入口
func main() {
	ctx := context.Background()

	// 开grpc 服务
	s := grpc.NewServer(
		//grpc.Creds(credentials.NewTLS(tlsCfg)),
		grpc.ConnectionTimeout(time.Second * 2),
	)
	grpcAddr := fmt.Sprintf("%s:%s",
		cfg.Get(ctx, cfg.GrpcHostKey, cfg.GrpcDefaultHost),
		cfg.Get(ctx, cfg.GrpcPortKey, strconv.Itoa(cfg.GrpcDefaultPort)),
	)

	// 注册服务
	proto.RegisterConfigServiceServer(s, &CfgService{})

	ls, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed listen %s %v", grpcAddr, err)
	}

	log.Printf("grpc server running on %s ...", grpcAddr)
	go func() {
		err = s.Serve(ls)
		if err != nil {
			log.Fatalf("failed serve grpc %v", err)
		}
	}()

	// grpc gateway 网关服务
	grpcGateway(ctx, grpcAddr)
}

func grpcGateway(ctx context.Context, grpcAddr string) {
	// grpc gate-way
	// 创建grpc-gateway网关http服务器
	mux := runtime.NewServeMux()
	// 获取grpc gateway地址
	gateWayAddr := fmt.Sprintf("%s:%s",
		cfg.Get(ctx, cfg.GrpcGateWayHostKey, cfg.GrpcGateWayDefaultHost),
		cfg.Get(ctx, cfg.GrpcGateWayPortKey, strconv.Itoa(cfg.GrpcGateWayDefaultPort)),
	)
	// 开服务
	ls, err := net.Listen("tcp", gateWayAddr)
	if err != nil {
		log.Fatalf("failed listen %s %v", gateWayAddr, err)
	}
	// 注册服务
	err = proto.RegisterConfigServiceHandlerFromEndpoint(ctx, mux, grpcAddr, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		log.Fatal("failed register grpc gateway server ", err)
	}
	err = http.Serve(ls, mux)
	if err != nil {
		log.Fatalf("failed server %s %v", gateWayAddr, err)
	}
}

// grpc 服务实现
type CfgService struct {
	proto.ConfigServiceServer
}

// 获取配置
func (*CfgService) GetConfig(ctx context.Context, req *proto.GetCfgReq) (*proto.GetCfgRsp, error) {
	return &proto.GetCfgRsp{
		Code:  req.Code,
		Value: cfg.Get(ctx, req.Code),
	}, nil
}

// 设置配置
func (*CfgService) SetConfig(ctx context.Context, req *proto.SetCfgReq) (*proto.SetCfgRsp, error) {
	err := cfg.Set(ctx, req.Code, req.Value)
	if err != nil {
		return nil, err
	}
	return &proto.SetCfgRsp{}, nil
}
