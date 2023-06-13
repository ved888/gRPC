package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc/pb_folder/pb"
	"log"
	"time"
)

// LaptopServer is the server that provides laptop service
type LaptopServer struct {
	Store LaptopStore
	pb.UnimplementedLaptopServiceServer
}

// NewLaptopServer return a new laptop server
func NewLaptopServer(store LaptopStore) *LaptopServer {
	return &LaptopServer{Store: store}
}

// CreateLaptop is an unary RPC to create a new laptop
func (server *LaptopServer) CreateLaptop(ctx context.Context, req *pb.CreateLaptopRequest) (*pb.CreateLaptopResponse, error) {
	laptop := req.GetLaptop()
	log.Printf("receive a laptop request with id :%s", laptop.Id)

	if len(laptop.Id) > 0 {
		// check if it's a valid uuid
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "laptop ID is not a valid UUID:%v", err)

		}

	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate a new laptop ID:%v", err)
		}
		laptop.Id = id.String()
	}
	//save to the laptop to the db.

	// some heavy processing
	time.Sleep(6 * time.Second)
	// when request canceled by client then server receive request but does not save laptop with id
	if ctx.Err() == context.Canceled {
		log.Print("request is canceled")
		return nil, status.Error(codes.Canceled, "request is canceled")
	}

	// when request canceled by the client then server is deadline but server receive request and save laptop with id
	if ctx.Err() == context.DeadlineExceeded {
		log.Print("deadline is exceeded")
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	//save to the laptop to the in-memory storage
	err := server.Store.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "cannot save laptop to the store:%v", err)
	}
	log.Printf("saved laptop with id:%s", laptop.Id)
	res := &pb.CreateLaptopResponse{
		Id: laptop.Id,
	}
	return res, nil

}
