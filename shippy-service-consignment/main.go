package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/micro/go-micro/v2"
	pb "github.com/raissov/shippy/shippy-service-consignment/proto/consignment"
	vesselProto "github.com/raissov/shippy/shippy-service-vessel/proto/vessel"
)

func main() {

	service := micro.NewService(
		micro.Name("shippy.service.consignment"),
	)

	service.Init()

	uri := os.Getenv("DB_HOST")

	client, err := CreateClient(context.Background(), uri, 0)
	if err != nil {
		log.Panic(err)
	}
	defer client.Disconnect(context.Background())

	consignmentCollection := client.Database("shippy").Collection("consignments")

	repository := &MongoRepository{consignmentCollection}
	vesselClient := vesselProto.NewVesselService("shippy.service.client", service.Client())
	h := &handler{repository, vesselClient}


	pb.RegisterShippingServiceHandler(service.Server(), h)


	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
