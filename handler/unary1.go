package handler

import (
	"context"
	"log"
	"time"
	pb "userclientservice/proto"
)

func CreateUser(grpcclient1 pb.Client1RequestClient, user *pb.UserDetails) (*pb.UserResponse1, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()
	res, err := grpcclient1.CreateUser(ctx, user)
	if err != nil {
		log.Fatalf("could not create the user: %v", err)
	}
	return res, nil
}
