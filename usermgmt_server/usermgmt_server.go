package main

import (
	"context"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"

	pb "github.com/blazingly-fast/user-management-service/usermgmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	port = ":50051"
)

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{
	}
}

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
}

func (server *UserManagementServer) Run() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserManagementServer(s, server)
    reflection.Register(s)
	log.Printf("server listening at %v", lis.Addr())
	return s.Serve(lis)
}

func (server *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Received: %v", in.GetName())
    readBytes, err := ioutil.ReadFile("users.json")
    var users_list *pb.UserList = &pb.UserList{}
	var user_id = int32(rand.Intn(100))
	created_user := &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: user_id}
   
    if err != nil{
    if os.IsNotExist(err) {
      log.Print("File not found. Creating a new file")
      users_list.Users = append(users_list.Users, created_user)
      jsonBytes, err := protojson.Marshal(users_list)
      if err != nil {
        log.Fatalf("JSON Marshaling failed: %v", err)
      }
      if err := ioutil.WriteFile("users.json", jsonBytes, 0664); err != nil {
        log.Fatalf("Failed write to file: %v", err)
      }
      return created_user, nil
      
    } else {
      log.Fatalln("Error reading file: ", err)
    }
  }
    if err = protojson.Unmarshal(readBytes, users_list); err != nil{
    log.Fatalf("Failed to parse user list: %v", err)
  }
  users_list.Users = append(users_list.Users, created_user)
  jsonBytes, err := protojson.Marshal(users_list)
  if err != nil {
    log.Fatalf("JSON Marshaling failed: %v", err)
  }
  if err := ioutil.WriteFile("users.json", jsonBytes, 0664); err != nil {
    log.Fatalf("Failed write to file: %v", err)
  }
	return created_user, nil
}

func (server *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
  jsonBytes, err := ioutil.ReadFile("users.json")
  if err != nil {
    log.Fatalf("Failed reading from file %v", err)
  }
  var users_list *pb.UserList = &pb.UserList{}
  if err := protojson.Unmarshal(jsonBytes, users_list); err != nil {
    log.Fatalf("Unmarshaling failed: %v", err)
  }

  return users_list, nil
}

func main() {
	var user_mgmt_server *UserManagementServer = NewUserManagementServer()
	if err := user_mgmt_server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
