package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/golang/protobuf/ptypes"

	dist "github.com/dm-alexi/go_piscine/rush00/ex02/distribution"
	pb "github.com/dm-alexi/go_piscine/rush00/ex02/transmitter"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

const (
	serverAddr   = "localhost:50051"
	initialTries = 100
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

// Anomaly struct in database
type Anomaly struct {
	SessionID string    `pg:"SessionID"`
	Frequency float64   `pg:"Frequency"`
	Time      time.Time `pg:"Time"`
}

func main() {
	var k float64
	flag.Float64Var(&k, "k", 5.0, "anomaly detection coefficient")
	flag.Parse()
	quants := make([]float64, initialTries)
	// connect to database
	db := pg.Connect(&pg.Options{
		Addr:     ":5432",
		User:     "rush",
		Password: "rush",
		Database: "rush",
	})
	defer db.Close()
	err := db.Model(&Anomaly{}).CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
	})
	if err != nil {
		log.Fatalf("Database error: %v", err)
	}
	// connect to emitter
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()
	client := pb.NewEmitterClient(conn)
	stream, err := client.BeginTransmission(context.Background(), &pb.ConnectionRequest{})
	if err != nil {
		log.Fatalf("failed to receive transmission: %v", err)
	}
	log.Println("Connection established, starting transmission")
	// collection stage
	for i := 0; i < initialTries; i++ {
		quant, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("%v error: %v", client, err)
		}
		quants[i] = quant.Frequency
	}
	nd := dist.AnalyzeRow(quants)
	log.Printf("Information gathered, stats guessed from %d quants: mean = %f, stddev = %f\n", initialTries, nd.Mean, nd.Stddev)
	// analyze stage
	for i := 0; ; i++ {
		quant, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("%v error: %v", client, err)
		}
		if i%(100*initialTries) == 0 {
			log.Printf("Processed signals: %d\n", i)
		}
		if quant.Frequency < nd.Mean-nd.Stddev*k || quant.Frequency > nd.Mean+nd.Stddev*k {
			log.Printf("Anomaly detected: %v", quant)
			_, err = db.Model(&Anomaly{SessionID: quant.GetSessionId(), Frequency: quant.GetFrequency(), Time: quant.GetTime().AsTime()}).Insert()
			if err != nil {
				log.Printf("Anomaly insertion error: %v", err)
			}
		}
	}
}
