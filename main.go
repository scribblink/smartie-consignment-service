package main

import (
	"fmt"
	"os"
	"log"

	// Import the generated protobuf code
	"context"

	"github.com/micro/go-micro"
	pb "github.com/scribblink/smartie-consignment-service/proto/consignment"
	vehicleProto "github.com/scribblink/smartie-vehicle-service/proto/vehicle"
)

const (
	defaultHost = "datastore:27017"
)

func main() {
	// Set-up micro instance
	srv := micro.NewService(
		micro.Name("smartie.consignment.service"),
	)

	srv.Init()

	uri := os.Getenv("DB_HOST")
	if uri == "" {
		uri = defaultHost
	}
	client, err := CreateClient(uri)
	if err != nil {
		log.Panic(err)
	}
	defer client.Disconnect(context.TODO())

	consignmentCollection := client.Database("smartie").Collection("consignments")

	repository := &MongoRepository{consignmentCollection}
	vehicleClient := vehicleProto.NewVehicleServiceClient("smartie.client.service", srv.Client())
	h := &handler{repository, vehicleClient}

	// Register handlers
	pb.RegisterShippingServiceHandler(srv.Server(), h)

	// Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
