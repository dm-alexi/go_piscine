package main

import (
	"log"
	"net"

	"github.com/golang/protobuf/ptypes"

	dist "github.com/dm-alexi/go_piscine/rush00/ex00/distribution"
	pb "github.com/dm-alexi/go_piscine/rush00/ex00/transmitter"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement transmitter.EmitterServer.
type server struct {
	pb.UnimplementedEmitterServer
	sessions map[uuid.UUID]*dist.NormDist
}

// BeginTransmission implements transmitter.BeginTransmission
func (s *server) BeginTransmission(in *pb.ConnectionRequest, stream pb.Emitter_BeginTransmissionServer) error {
	id := uuid.New()
	s.sessions[id] = dist.GetDistribution()
	log.Printf("Established new session: id = %v, mean = %v, stddev = %v", id, s.sessions[id].Mean, s.sessions[id].Stddev)
	for {
		entry := pb.Quant{
			SessionId: id.String(),
			Frequency: dist.GetEntry(s.sessions[id]),
			Time:      ptypes.TimestampNow()}
		if err := stream.Send(&entry); err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterEmitterServer(s, &server{sessions: make(map[uuid.UUID]*dist.NormDist)})
	log.Printf("Server started, listening on %v", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
