package main

import (
	"context"
	pb "github.com/scribblink/smartie-consignment-service/proto/consignment"
	vehicleProto "github.com/scribblink/smartie-vehicle-service/proto/vehicle"
	"log"
)

type handler struct {
	Repository
	vehicleClient vehicleProto.VehicleServiceClient
}

// CreateConsignment - we created just one method on our service,
// which is a create method, which takes a context and a request as an
// argument, these are handled by the gRPC server.
func (s *handler) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {

	// Here we call a client instance of our vehicle service with our consignment weight,
	// and the amount of containers as the capacity value
	vehicleResponse, err := s.vehicleClient.FindAvailable(ctx, &vehicleProto.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})
	log.Printf("Found vehicle: %s \n", vehicleResponse.Vehicle.Name)
	if err != nil {
		return err
	}

	// We set the VehicleId as the vehicle we got back from our
	// vehicle service
	req.VehicleId = vehicleResponse.Vehicle.Id

	// Save our consignment
	if err = s.Repository.Create(req); err != nil {
		return err
	}

	res.Created = true
	res.Consignment = req
	return nil
}

// GetConsignments -
func (s *handler) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments, err := s.Repository.GetAll()
	if err != nil {
		return err
	}
	res.Consignments = consignments
	return nil
}
