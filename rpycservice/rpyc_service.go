package rpycservice

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type XRayService struct {
	core       *XRayCore
	connection *grpc.ClientConn
	UnimplementedXRayServiceServer
}

func NewXRayService() *XRayService {
	return &XRayService{}
}

func (s *XRayService) OnConnect(conn *grpc.ClientConn) {
	s.connection = conn
	// Additional logic if needed
}

func (s *XRayService) OnDisconnect() {
	// Additional logic if needed
	if s.core != nil {
		s.core.Stop()
	}
	s.core = nil
	s.connection = nil
}

func (s *XRayService) Start(ctx context.Context, req *StartRequest) (*StartResponse, error) {
	config := req.Config
	// Convert config and start XRayCore
	log.Println("Starting with config:", config)
	// Additional start logic
	return &StartResponse{Message: "Started"}, nil
}

func (s *XRayService) Stop(ctx context.Context, req *StopRequest) (*StopResponse, error) {
	if s.core != nil {
		s.core.Stop()
	}
	return &StopResponse{Message: "Stopped"}, nil
}

func (s *XRayService) Restart(ctx context.Context, req *RestartRequest) (*RestartResponse, error) {
	config := req.Config
	// Convert config and restart XRayCore
	log.Println("Restarting with config:", config)
	// Additional restart logic
	return &RestartResponse{Message: "Restarted"}, nil
}

func (s *XRayService) FetchXRayVersion(ctx context.Context, req *FetchXRayVersionRequest) (*FetchXRayVersionResponse, error) {
	if s.core == nil {
		return nil, fmt.Errorf("Xray has not been started")
	}
	return &FetchXRayVersionResponse{Version: s.core.Version()}, nil
}

func (s *XRayService) FetchLogs(ctx context.Context, req *FetchLogsRequest) (*FetchLogsResponse, error) {
	if s.core != nil {
		handler := NewXRayCoreLogsHandler(s.core, func(log string) {
			// Send logs back to the client
		}, time.Second*1)
		defer handler.Stop()
	}
	return &FetchLogsResponse{Message: "Logs streaming"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	RegisterXRayServiceServer(grpcServer, NewXRayService())
	log.Println("Server listening on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
