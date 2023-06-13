package service

import (
	"bytes"
	"context"
	"errors"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc/pb/pb"
	"io"
	"log"
)

// maximum 1 megabyte
const maxImageSize = 1 << 20

// LaptopServer is the server that provides laptop service
type LaptopServer struct {
	laptopStore LaptopStore
	imageStore  ImageStore
	pb.UnimplementedLaptopServiceServer
}

// NewLaptopServer return a new laptop server
func NewLaptopServer(laptopStore LaptopStore, imageStore ImageStore) *LaptopServer {
	return &LaptopServer{laptopStore: laptopStore, imageStore: imageStore}
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
	//	time.Sleep(6 * time.Second)

	if err := contextError(ctx); err != nil {
		return nil, err
	}

	//save to the laptop to the in-memory storage
	err := server.laptopStore.Save(laptop)
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

func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return logError(status.Error(codes.Canceled, "request is canceled"))
	// when request canceled by the client then server is deadline but server receive request and save laptop with id
	case context.DeadlineExceeded:
		return logError(status.Error(codes.DeadlineExceeded, "deadline is exceeded"))
	default:
		return nil
	}

}

// SearchLaptop is a server-streaming RPC to search for laptop
func (server *LaptopServer) SearchLaptop(req *pb.SearchLaptopRequest, stream pb.LaptopService_SearchLaptopServer) error {
	filter := req.GetFilter()
	log.Printf("receive a search-laptop request with filter:%v", filter)

	err := server.laptopStore.Search(
		stream.Context(),
		filter,
		func(laptop *pb.Laptop) error {
			res := &pb.SearchLaptopResponse{Laptop: laptop}

			err := stream.Send(res)
			if err != nil {
				return err
			}
			log.Printf("sent laptop with id %v", laptop.GetId())
			return nil
		},
	)
	if err != nil {
		return status.Errorf(codes.Internal, "uexpected error %v", err)

	}
	return nil
}

// UploadImage is a client-streaming to upload a laptop image
func (server *LaptopServer) UploadImage(stream pb.LaptopService_UploadImageServer) error {
	req, err := stream.Recv()
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot receive image info"))
	}

	laptopId := req.GetInfo().GetLaptopId()
	imageType := req.GetInfo().GetImageType()
	log.Printf("receive an upload image request for laptop %s with image type %s", laptopId, imageType)

	laptop, err := server.laptopStore.Find(laptopId)
	if err != nil {
		return logError(status.Errorf(codes.Internal, "cannot find laptop:%v", err))
	}
	if laptop == nil {
		return logError(status.Errorf(codes.InvalidArgument, "laptop %s does't exist", laptopId))
	}

	imageData := bytes.Buffer{}
	imageSize := 0

	for {
		// check context error
		if err := contextError(stream.Context()); err != nil {
			return nil
		}
		log.Printf("waiting to reacive more data")
		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("no more data")
			break
		}
		if err != nil {
			return logError(status.Errorf(codes.Unknown, "cannot receive chunk data:%v", err))

		}
		chunk := req.GetChunkData()
		size := len(chunk)

		log.Printf("received a chunk with size %d", size)

		imageSize += size
		if imageSize > maxImageSize {
			return logError(status.Errorf(codes.InvalidArgument, "image is too larg:%d>%d", imageSize, maxImageSize))
		}
		// write slowly
		//		time.Sleep(time.Second)

		_, err = imageData.Write(chunk)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "cannot write chunk data:%v", err))
		}
	}
	imageId, err := server.imageStore.Save(laptopId, imageType, imageData)
	if err != nil {
		return logError(status.Errorf(codes.Internal, "cannot save image to the store:%v", err))
	}
	res := &pb.UploadImageResponse{
		Id:   imageId,
		Size: uint32(imageSize),
	}
	err = stream.SendAndClose(res)
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot send response:%v", err))
	}
	log.Printf("saved image with id:%s,size:%d", imageId, imageSize)
	return nil

}

func logError(err error) error {
	if err != nil {
		log.Print(err)
	}
	return err
}
