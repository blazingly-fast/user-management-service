package main

import (
	"context"
	"log"
	"math/rand"
	"net"

	pb "github.com/blazingly-fast/user-management-service/usermgmt"
	"google.golang.org/grpc"
)



const (
port = ":50051"
)

type UserManagementServer struct {
 pb.UnimplementedUserManagementServer 
}

func (u*UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
  log.Printf("Recived: %v", in.GetName())
  user_id := rand.Int31n(1000)
  return &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: user_id}, nil
}

func main(){
  lis, err := net.Listen("tcp", port)
  if err != nil {
    log.Fatalf("failed to listen %v", err)
  }

  s:= grpc.NewServer()
  pb.RegisterUserManagementServer(s, &UserManagementServer{})
  log.Printf("server listening at %v", lis.Addr())
  if err := s.Serve(lis); err != nil {
    log.Fatalf("falied to serve: %v", err)
  }
}
