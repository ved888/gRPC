package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"grpc/pb_folder/pb"
	"grpc/sample"
	"io"
	"log"
	"time"
)

func createLaptop(laptopClient pb.LaptopServiceClient) {
	laptop := sample.NewLaptop()
	// put the laptop id and you can check from here laptop exist or not
	// if not give empty string so when run client then auto created laptop id
	// if you give any new uuid then it is created laptop id by given uuid
	// we give request to server
	laptop.Id = ""
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	// set timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := laptopClient.CreateLaptop(ctx, req)
	//res, err := laptopClient.CreateLaptop(context.Background(), req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			// not a big deal
			log.Printf("laptop already exists")
		} else {
			log.Fatal("cannot create server:", err)

		}
		return
	}
	log.Printf("create laptop with id:%s", res.Id)
}

func searchLaptop(laptopClient pb.LaptopServiceClient, filter *pb.Filter) {
	log.Print("search filter", filter)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.SearchLaptopRequest{Filter: filter}
	stream, err := laptopClient.SearchLaptop(ctx, req)
	if err != nil {
		log.Fatal("cannot search laptop:", err)
	}
	for {
		res, err := stream.Recv()
		if err != io.EOF {
			return
		}
		if err != nil {
			log.Fatal("cannot receive response", err)
		}
		laptop := res.GetLaptop()
		log.Print("found:", laptop.GetId())
		log.Print("brand", laptop.GetBrand())
		log.Print("name", laptop.GetName())
		log.Print("cpu cores", laptop.GetCpu().GetNumberCores())
		log.Print("cpu min ghz", laptop.GetCpu().GetMinGhz())
		log.Print("ram", laptop.GetRam().GetValue(), laptop.GetRam().GetUnit())
		log.Print("price", laptop.GetPriceUsd(), "usd")
	}
}

func main() {
	// Connection to internal grpc server
	conn1, err := grpc.Dial("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	laptopClient := pb.NewLaptopServiceClient(conn1)
	// if make file work then above coon1 remove from hear  and below commented line do uncomment

	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dail server %s", *serverAddress)

	//	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server:", err)
	}
	//	laptopClient := pb.NewLaptopServiceClient(conn)

	for i := 0; i < 10; i++ {
		createLaptop(laptopClient)
	}
	filter := &pb.Filter{
		MaxPriceUsd: 3000,
		MinCpuCores: 4,
		MinCpuGhz:   2.5,
		MinRam:      &pb.Memory{Value: 8, Unit: pb.Memory_GIGABYTE},
	}
	searchLaptop(laptopClient, filter)
}
