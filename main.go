package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"userclientservice/handler"
	pb "userclientservice/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	port = ":8083"
)

type Request struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GrpcConnection() (pb.Client1RequestClient, *grpc.ClientConn) {
	conn, err := grpc.Dial("localhost"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Connection failed %v", err)
	}

	return pb.NewClient1RequestClient(conn), conn
}

func main() {
	client1, conn := GrpcConnection()
	defer conn.Close()
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var req Request
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userData := &pb.UserDetails{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		}

		res, err := handler.CreateUser(client1, userData)
		if err != nil {
			log.Printf("cannot create the user %v", err)
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(res)
	})

	fmt.Println("server running at 8030")

	if err := http.ListenAndServe(":8030", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	fmt.Println("its working")
}
