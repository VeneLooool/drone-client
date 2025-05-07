package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/VeneLooool/drone-client/internal/app/api/v1/drones"
	"github.com/VeneLooool/drone-client/internal/config"
	"github.com/VeneLooool/drone-client/internal/kafka/external-drone-events/publisher"
	drones_pb "github.com/VeneLooool/drone-client/internal/pb/api/v1/drones"
	drone_uc "github.com/VeneLooool/drone-client/internal/usecase/drone"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.New(ctx)
	if err != nil {
		log.Fatalf("failed to create new config: %s", err.Error())
	}

	go func() {
		if err := runGRPC(ctx, cfg); err != nil {
			log.Fatal(err)
		}
	}()

	if err := runHTTPGateway(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}

func runGRPC(ctx context.Context, cfg *config.Config) error {
	grpcServer := grpc.NewServer()
	defer grpcServer.GracefulStop()

	droneServer, err := newServices(ctx, cfg)
	if err != nil {
		return err
	}
	drones_pb.RegisterDronesServer(grpcServer, droneServer)

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %s", err.Error())
	}

	log.Printf("gRPC server listening on :%s\n", cfg.GrpcPort)
	if err = grpcServer.Serve(grpcListener); err != nil {
		return err
	}
	return nil
}

func runHTTPGateway(ctx context.Context, cfg *config.Config) error {
	mux := runtime.NewServeMux()
	err := drones_pb.RegisterDronesHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%s", cfg.GrpcPort), []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})
	if err != nil {
		log.Fatalf("failed to register gateway: %s", err.Error())
	}

	// Serve Swagger JSON and Swagger UI
	fs := http.FileServer(http.Dir("./swagger-ui")) // директория со статикой UI
	http.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", fs))

	// Serve Swagger JSON файл
	http.HandleFunc("/swagger/drones.swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./internal/pb/api/v1/drones/drones.swagger.json")
	})

	withCORS := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			// Для preflight-запросов
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			h.ServeHTTP(w, r)
		})
	}

	// gRPC → REST mux
	http.Handle("/", withCORS(mux))

	log.Printf("HTTP gateway listening on :%s\n", cfg.HttpPort)
	if err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.HttpPort), nil); err != nil {
		return err
	}

	return nil
}

func newServices(ctx context.Context, cfg *config.Config) (*drones.Implementation, error) {
	droneEventsPublisher := publisher.New(ctx, cfg.GetKafkaConfig())

	droneUC := drone_uc.New(droneEventsPublisher)

	return drones.NewService(droneUC), nil
}
